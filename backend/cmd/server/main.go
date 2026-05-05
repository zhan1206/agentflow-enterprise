// Package main is the entry point for AgentFlow-Enterprise API server
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zhan1206/agentflow-enterprise/backend/internal/core/scheduler"
	"github.com/zhan1206/agentflow-enterprise/backend/internal/core/sync"
	"github.com/zhan1206/agentflow-enterprise/backend/internal/observability/tracing"
	"github.com/zhan1206/agentflow-enterprise/backend/internal/security/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	configPath = flag.String("config", "config.yaml", "Path to configuration file")
	version    = "0.1.0-dev"
)

func main() {
	flag.Parse()

	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting AgentFlow-Enterprise server",
		zap.String("version", version),
		zap.String("config", *configPath),
	)

	// Initialize tracing
	shutdownTracing, err := tracing.InitTracer("agentflow-server")
	if err != nil {
		logger.Fatal("Failed to initialize tracing", zap.Error(err))
	}
	defer shutdownTracing()

	// Initialize state sync engine (Raft-based)
	syncEngine, err := sync.NewEngine(sync.Config{
		NodeID:    "node-1",
		DataDir:   "./data/raft",
		Bootstrap: true,
	})
	if err != nil {
		logger.Fatal("Failed to initialize sync engine", zap.Error(err))
	}
	defer syncEngine.Shutdown()

	// Initialize task scheduler
	taskScheduler := scheduler.New(scheduler.Config{
		Workers:     10,
		QueueSize:   1000,
		SyncEngine:  syncEngine,
		Logger:      logger,
	})
	defer taskScheduler.Shutdown()

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(auth.Config{
		JWTSecret: "dev-secret-change-in-production",
		Logger:    logger,
	})

	// Setup Gin router
	router := setupRouter(logger, authMiddleware, taskScheduler, syncEngine)

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		logger.Info("HTTP server listening", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}

	logger.Info("Server stopped")
}

func setupRouter(
	logger *zap.Logger,
	authMiddleware *auth.Middleware,
	scheduler *scheduler.Scheduler,
	syncEngine *sync.Engine,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(requestLogger(logger))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": version,
		})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// Public endpoints
		v1.POST("/auth/login", authMiddleware.LoginHandler)
		v1.POST("/auth/register", authMiddleware.RegisterHandler)

		// Protected endpoints
		protected := v1.Group("")
		protected.Use(authMiddleware.Authenticate())
		{
			// Agents
			agents := protected.Group("/agents")
			{
				agents.GET("", listAgents)
				agents.POST("", registerAgent)
				agents.GET("/:id", getAgent)
				agents.PUT("/:id", updateAgent)
				agents.DELETE("/:id", deleteAgent)
			}

			// Tasks
			tasks := protected.Group("/tasks")
			{
				tasks.GET("", listTasks)
				tasks.POST("", createTask)
				tasks.GET("/:id", getTask)
				tasks.POST("/:id/execute", executeTask(scheduler))
				tasks.POST("/:id/cancel", cancelTask(scheduler))
			}

			// Workflows
			workflows := protected.Group("/workflows")
			{
				workflows.GET("", listWorkflows)
				workflows.POST("", createWorkflow)
				workflows.GET("/:id", getWorkflow)
				workflows.PUT("/:id", updateWorkflow)
				workflows.DELETE("/:id", deleteWorkflow)
			}

			// Metrics & Observability
			observability := protected.Group("/observability")
			{
				observability.GET("/metrics", getMetrics)
				observability.GET("/traces/:taskId", getTraces)
				observability.GET("/dashboard", getDashboard)
			}
		}
	}

	return router
}

func requestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		logger.Info("HTTP request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
		)
	}
}

// Handler stubs (to be implemented)
func listAgents(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"agents": []string{}}) }
func registerAgent(c *gin.Context) { c.JSON(http.StatusCreated, gin.H{"id": "agent-1"}) }
func getAgent(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"id": c.Param("id")}) }
func updateAgent(c *gin.Context)   { c.JSON(http.StatusOK, gin.H{"id": c.Param("id")}) }
func deleteAgent(c *gin.Context)   { c.JSON(http.StatusNoContent, nil) }

func listTasks(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"tasks": []string{}}) }
func createTask(c *gin.Context)   { c.JSON(http.StatusCreated, gin.H{"id": "task-1"}) }
func getTask(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"id": c.Param("id")}) }
func executeTask(s *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{"status": "executing"})
	}
}
func cancelTask(s *scheduler.Scheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
	}
}

func listWorkflows(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"workflows": []string{}}) }
func createWorkflow(c *gin.Context)   { c.JSON(http.StatusCreated, gin.H{"id": "workflow-1"}) }
func getWorkflow(c *gin.Context)      { c.JSON(http.StatusOK, gin.H{"id": c.Param("id")}) }
func updateWorkflow(c *gin.Context)   { c.JSON(http.StatusOK, gin.H{"id": c.Param("id")}) }
func deleteWorkflow(c *gin.Context)   { c.JSON(http.StatusNoContent, nil) }

func getMetrics(c *gin.Context)   { c.JSON(http.StatusOK, gin.H{"metrics": {}}) }
func getTraces(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"traces": []string{}}) }
func getDashboard(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"dashboard": {}}) }
