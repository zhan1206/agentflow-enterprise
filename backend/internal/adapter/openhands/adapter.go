// Package openhands provides an adapter for the OpenHands framework
package openhands

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/agentflow-enterprise/agentflow-enterprise/backend/internal/adapter/base"
)

// Adapter implements the OpenHands adapter
type Adapter struct {
	*base.BaseAdapter
	endpoint   string
	httpClient *http.Client
	apiKey     string
}

// Config holds OpenHands adapter configuration
type Config struct {
	Endpoint string
	APIKey   string
	Timeout  time.Duration
}

// New creates a new OpenHands adapter
func New(cfg Config) *Adapter {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	return &Adapter{
		BaseAdapter: base.NewBaseAdapter(base.AdapterInfo{
			Name:    "OpenHands Adapter",
			Type:    "openhands",
			Version: "1.0.0",
			Capabilities: []string{
				"code-generation",
				"code-review",
				"debugging",
				"testing",
				"documentation",
			},
		}),
		endpoint: cfg.Endpoint,
		apiKey:   cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// Execute executes a task on an OpenHands agent
func (a *Adapter) Execute(ctx context.Context, agentID string, task base.TaskRequest) (*base.TaskResult, error) {
	startTime := time.Now()

	// Prepare request
	reqBody := map[string]interface{}{
		"task_id":     task.ID,
		"type":        task.Type,
		"description": task.Description,
		"input":       task.Input,
	}

	// Make HTTP request to OpenHands
	url := fmt.Sprintf("%s/api/v1/agents/%s/execute", a.endpoint, agentID)

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Status     string                 `json:"status"`
		Output     map[string]interface{} `json:"output"`
		Error      string                 `json:"error,omitempty"`
		TokenUsage int64                  `json:"token_usage"`
		Artifacts  []base.Artifact        `json:"artifacts,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &base.TaskResult{
		ID:         task.ID,
		Status:     base.ResultStatus(result.Status),
		Output:     result.Output,
		Error:      result.Error,
		Duration:   time.Since(startTime),
		TokenUsage: result.TokenUsage,
		Artifacts:  result.Artifacts,
	}, nil
}

// Cancel cancels a running task
func (a *Adapter) Cancel(ctx context.Context, taskID string) error {
	url := fmt.Sprintf("%s/api/v1/tasks/%s/cancel", a.endpoint, taskID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cancel request failed with status %d", resp.StatusCode)
	}

	return nil
}

// Health checks the health of the OpenHands service
func (a *Adapter) Health(ctx context.Context) (*base.HealthStatus, error) {
	url := fmt.Sprintf("%s/health", a.endpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return &base.HealthStatus{
			Healthy: false,
			Message: fmt.Sprintf("connection failed: %v", err),
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &base.HealthStatus{
			Healthy: false,
			Message: fmt.Sprintf("service returned status %d", resp.StatusCode),
		}, nil
	}

	return &base.HealthStatus{
		Healthy: true,
		Message: "OpenHands service is healthy",
	}, nil
}

// Stream streams task output
func (a *Adapter) Stream(ctx context.Context, taskID string) (<-chan base.StreamEvent, error) {
	events := make(chan base.StreamEvent, 100)

	url := fmt.Sprintf("%s/api/v1/tasks/%s/stream", a.endpoint, taskID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Start goroutine to read SSE stream
	go func() {
		defer close(events)
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		for {
			var event base.StreamEvent
			if err := decoder.Decode(&event); err != nil {
				if err == io.EOF {
					break
				}
				continue
			}
			event.Timestamp = time.Now()
			events <- event
		}
	}()

	return events, nil
}

// Close closes the adapter
func (a *Adapter) Close() error {
	a.httpClient.CloseIdleConnections()
	return nil
}

// Note: bytes package import needed
import "bytes"
