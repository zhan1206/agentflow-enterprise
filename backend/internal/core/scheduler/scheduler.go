// Package scheduler implements the distributed task scheduling engine
package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/zhan1206/agentflow-enterprise/backend/internal/core/sync"
	"go.uber.org/zap"
)

// Config holds scheduler configuration
type Config struct {
	Workers    int
	QueueSize  int
	SyncEngine *sync.Engine
	Logger     *zap.Logger
}

// Task represents a schedulable task
type Task struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"` // serial, parallel, dag
	Priority    int                    `json:"priority"`
	Agents      []string               `json:"agents"`
	Input       map[string]interface{} `json:"input"`
	Output      map[string]interface{} `json:"output,omitempty"`
	Status      TaskStatus             `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	StartedAt   *time.Time             `json:"startedAt,omitempty"`
	CompletedAt *time.Time             `json:"completedAt,omitempty"`
	Error       string                 `json:"error,omitempty"`
}

// TaskStatus represents the current status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusQueued    TaskStatus = "queued"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCancelled TaskStatus = "cancelled"
)

// SubTask represents a unit of work within a task
type SubTask struct {
	ID        string                 `json:"id"`
	TaskID    string                 `json:"taskId"`
	AgentID   string                 `json:"agentId"`
	Name      string                 `json:"name"`
	Input     map[string]interface{} `json:"input"`
	Output    map[string]interface{} `json:"output,omitempty"`
	Status    TaskStatus             `json:"status"`
	DependsOn []string               `json:"dependsOn,omitempty"`
}

// Scheduler is the main task scheduling engine
type Scheduler struct {
	config     Config
	logger     *zap.Logger
	syncEngine *sync.Engine

	taskQueue chan *Task
	workers   int

	mu          sync.RWMutex
	tasks       map[string]*Task
	subTasks    map[string]*SubTask
	taskResults map[string]interface{}

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// New creates a new scheduler instance
func New(config Config) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		config:      config,
		logger:      config.Logger,
		syncEngine:  config.SyncEngine,
		taskQueue:   make(chan *Task, config.QueueSize),
		workers:     config.Workers,
		tasks:       make(map[string]*Task),
		subTasks:    make(map[string]*SubTask),
		taskResults: make(map[string]interface{}),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start begins the scheduler workers
func (s *Scheduler) Start() {
	s.logger.Info("Starting scheduler", zap.Int("workers", s.workers))

	for i := 0; i < s.workers; i++ {
		s.wg.Add(1)
		go s.worker(i)
	}

	s.logger.Info("Scheduler started successfully")
}

// Shutdown gracefully stops the scheduler
func (s *Scheduler) Shutdown() {
	s.logger.Info("Shutting down scheduler...")
	s.cancel()
	s.wg.Wait()
	s.logger.Info("Scheduler shutdown complete")
}

// Submit adds a new task to the queue
func (s *Scheduler) Submit(task *Task) error {
	task.Status = TaskStatusQueued
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	s.mu.Lock()
	s.tasks[task.ID] = task
	s.mu.Unlock()

	// Persist to sync engine
	if s.syncEngine != nil {
		if err := s.syncEngine.Set("task:"+task.ID, task); err != nil {
			s.logger.Error("Failed to persist task", zap.Error(err), zap.String("taskId", task.ID))
		}
	}

	select {
	case s.taskQueue <- task:
		s.logger.Info("Task submitted", zap.String("taskId", task.ID))
		return nil
	default:
		return ErrQueueFull
	}
}

// Cancel cancels a running task
func (s *Scheduler) Cancel(taskID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return ErrTaskNotFound
	}

	if task.Status == TaskStatusCompleted {
		return ErrTaskAlreadyCompleted
	}

	task.Status = TaskStatusCancelled
	now := time.Now()
	task.CompletedAt = &now
	task.UpdatedAt = now

	s.logger.Info("Task cancelled", zap.String("taskId", taskID))
	return nil
}

// Get retrieves a task by ID
func (s *Scheduler) Get(taskID string) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// List returns all tasks
func (s *Scheduler) List() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (s *Scheduler) worker(id int) {
	defer s.wg.Done()

	logger := s.logger.With(zap.Int("workerId", id))
	logger.Info("Worker started")

	for {
		select {
		case <-s.ctx.Done():
			logger.Info("Worker shutting down")
			return
		case task := <-s.taskQueue:
			s.processTask(task, logger)
		}
	}
}

func (s *Scheduler) processTask(task *Task, logger *zap.Logger) {
	logger.Info("Processing task", zap.String("taskId", task.ID))

	s.mu.Lock()
	task.Status = TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now
	task.UpdatedAt = now
	s.mu.Unlock()

	// Simulate task processing based on type
	switch task.Type {
	case "serial":
		s.processSerial(task, logger)
	case "parallel":
		s.processParallel(task, logger)
	case "dag":
		s.processDAG(task, logger)
	default:
		s.processSerial(task, logger)
	}
}

func (s *Scheduler) processSerial(task *Task, logger *zap.Logger) {
	for _, agentID := range task.Agents {
		select {
		case <-s.ctx.Done():
			return
		default:
			s.processWithAgent(task, agentID, logger)
		}
	}
	s.completeTask(task, nil)
}

func (s *Scheduler) processParallel(task *Task, logger *zap.Logger) {
	var wg sync.WaitGroup
	for _, agentID := range task.Agents {
		wg.Add(1)
		go func(aid string) {
			defer wg.Done()
			s.processWithAgent(task, aid, logger)
		}(agentID)
	}
	wg.Wait()
	s.completeTask(task, nil)
}

func (s *Scheduler) processDAG(task *Task, logger *zap.Logger) {
	// DAG processing with topological sort
	// This is a simplified implementation
	s.processSerial(task, logger)
}

func (s *Scheduler) processWithAgent(task *Task, agentID string, logger *zap.Logger) {
	logger.Info("Processing with agent",
		zap.String("taskId", task.ID),
		zap.String("agentId", agentID),
	)

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
}

func (s *Scheduler) completeTask(task *Task, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	task.CompletedAt = &now
	task.UpdatedAt = now

	if err != nil {
		task.Status = TaskStatusFailed
		task.Error = err.Error()
	} else {
		task.Status = TaskStatusCompleted
	}

	s.logger.Info("Task completed",
		zap.String("taskId", task.ID),
		zap.String("status", string(task.Status)),
	)

	// Persist final state
	if s.syncEngine != nil {
		s.syncEngine.Set("task:"+task.ID, task)
	}
}

// Errors
var (
	ErrQueueFull             = &SchedulerError{Message: "task queue is full"}
	ErrTaskNotFound          = &SchedulerError{Message: "task not found"}
	ErrTaskAlreadyCompleted  = &SchedulerError{Message: "task already completed"}
	ErrInvalidTaskType       = &SchedulerError{Message: "invalid task type"}
	ErrCircularDependency    = &SchedulerError{Message: "circular dependency detected"}
)

// SchedulerError represents a scheduler error
type SchedulerError struct {
	Message string
}

func (e *SchedulerError) Error() string {
	return e.Message
}
