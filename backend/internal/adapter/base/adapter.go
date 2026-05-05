// Package base provides the base adapter interface for Agent frameworks
package base

import (
	"context"
	"time"
)

// AgentStatus represents the current status of an Agent
type AgentStatus string

const (
	StatusOnline  AgentStatus = "online"
	StatusOffline AgentStatus = "offline"
	StatusBusy    AgentStatus = "busy"
	StatusError   AgentStatus = "error"
)

// AgentInfo contains information about an Agent
type AgentInfo struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`        // openhands, langgraph, crewai, autogen, custom
	Version      string            `json:"version"`
	Status       AgentStatus       `json:"status"`
	Capabilities []string          `json:"capabilities"` // skills, tools, etc.
	Metadata     map[string]string `json:"metadata"`
	LastSeen     time.Time         `json:"lastSeen"`
	Endpoint     string            `json:"endpoint"`
}

// TaskRequest represents a request to execute a task
type TaskRequest struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Input       map[string]interface{} `json:"input"`
	Timeout     time.Duration          `json:"timeout"`
	Context     context.Context        `json:"-"`
}

// TaskResult represents the result of a task execution
type TaskResult struct {
	ID          string                 `json:"id"`
	Status      ResultStatus           `json:"status"`
	Output      map[string]interface{} `json:"output"`
	Error       string                 `json:"error,omitempty"`
	Duration    time.Duration          `json:"duration"`
	TokenUsage  int64                  `json:"tokenUsage"`
	Artifacts   []Artifact             `json:"artifacts,omitempty"`
}

// ResultStatus represents the status of a task result
type ResultStatus string

const (
	ResultSuccess ResultStatus = "success"
	ResultFailure ResultStatus = "failure"
	ResultTimeout ResultStatus = "timeout"
	ResultPartial ResultStatus = "partial"
)

// Artifact represents a generated artifact
type Artifact struct {
	Type        string `json:"type"` // file, url, data
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	URL         string `json:"url,omitempty"`
	Content     string `json:"content,omitempty"`
}

// HealthStatus represents the health check result
type HealthStatus struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Adapter is the interface that all Agent framework adapters must implement
type Adapter interface {
	// Info returns information about this adapter
	Info() AdapterInfo

	// Register registers an Agent with the adapter
	Register(ctx context.Context, agent AgentInfo) error

	// Unregister removes an Agent from the adapter
	Unregister(ctx context.Context, agentID string) error

	// GetAgent retrieves an Agent by ID
	GetAgent(ctx context.Context, agentID string) (*AgentInfo, error)

	// ListAgents lists all registered Agents
	ListAgents(ctx context.Context) ([]AgentInfo, error)

	// Execute executes a task on an Agent
	Execute(ctx context.Context, agentID string, task TaskRequest) (*TaskResult, error)

	// Cancel cancels a running task
	Cancel(ctx context.Context, taskID string) error

	// Health checks the health of the adapter
	Health(ctx context.Context) (*HealthStatus, error)

	// Stream streams task output (optional)
	Stream(ctx context.Context, taskID string) (<-chan StreamEvent, error)

	// Close closes the adapter
	Close() error
}

// AdapterInfo contains information about an adapter
type AdapterInfo struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Version      string   `json:"version"`
	Capabilities []string `json:"capabilities"`
}

// StreamEvent represents a streaming event
type StreamEvent struct {
	Type      string      `json:"type"` // output, error, complete
	TaskID    string      `json:"taskId"`
	Data      string      `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// BaseAdapter provides common functionality for adapters
type BaseAdapter struct {
	info    AdapterInfo
	agents  map[string]*AgentInfo
}

// NewBaseAdapter creates a new BaseAdapter
func NewBaseAdapter(info AdapterInfo) *BaseAdapter {
	return &BaseAdapter{
		info:   info,
		agents: make(map[string]*AgentInfo),
	}
}

// Info returns adapter information
func (a *BaseAdapter) Info() AdapterInfo {
	return a.info
}

// Register registers an Agent
func (a *BaseAdapter) Register(ctx context.Context, agent AgentInfo) error {
	agent.LastSeen = time.Now()
	agent.Status = StatusOnline
	a.agents[agent.ID] = &agent
	return nil
}

// Unregister removes an Agent
func (a *BaseAdapter) Unregister(ctx context.Context, agentID string) error {
	delete(a.agents, agentID)
	return nil
}

// GetAgent retrieves an Agent
func (a *BaseAdapter) GetAgent(ctx context.Context, agentID string) (*AgentInfo, error) {
	agent, exists := a.agents[agentID]
	if !exists {
		return nil, ErrAgentNotFound
	}
	return agent, nil
}

// ListAgents lists all Agents
func (a *BaseAdapter) ListAgents(ctx context.Context) ([]AgentInfo, error) {
	result := make([]AgentInfo, 0, len(a.agents))
	for _, agent := range a.agents {
		result = append(result, *agent)
	}
	return result, nil
}

// Errors
var (
	ErrAgentNotFound    = &AdapterError{Code: "AGENT_NOT_FOUND", Message: "agent not found"}
	ErrTaskNotFound     = &AdapterError{Code: "TASK_NOT_FOUND", Message: "task not found"}
	ErrAgentBusy        = &AdapterError{Code: "AGENT_BUSY", Message: "agent is busy"}
	ErrTimeout          = &AdapterError{Code: "TIMEOUT", Message: "operation timed out"}
	ErrInvalidRequest   = &AdapterError{Code: "INVALID_REQUEST", Message: "invalid request"}
	ErrInternalError    = &AdapterError{Code: "INTERNAL_ERROR", Message: "internal error"}
)

// AdapterError represents an adapter error
type AdapterError struct {
	Code    string
	Message string
}

func (e *AdapterError) Error() string {
	return e.Message
}
