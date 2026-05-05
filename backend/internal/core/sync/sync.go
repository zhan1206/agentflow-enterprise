// Package sync implements distributed state synchronization using Raft
package sync

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Config holds sync engine configuration
type Config struct {
	NodeID    string
	DataDir   string
	Bootstrap bool
	Peers     []string
	Logger    *zap.Logger
}

// State represents a synchronized state entry
type State struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	Version   int64       `json:"version"`
	UpdatedAt time.Time   `json:"updatedAt"`
	UpdatedBy string      `json:"updatedBy"`
}

// Engine is the distributed state synchronization engine
type Engine struct {
	config Config
	logger *zap.Logger

	mu     sync.RWMutex
	state  map[string]*State
	watchers map[string][]Watcher

	ctx    context.Context
	cancel context.CancelFunc
}

// Watcher is a callback for state changes
type Watcher func(key string, value interface{}, version int64)

// NewEngine creates a new sync engine
func NewEngine(config Config) (*Engine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	e := &Engine{
		config:   config,
		logger:   config.Logger,
		state:    make(map[string]*State),
		watchers: make(map[string][]Watcher),
		ctx:      ctx,
		cancel:   cancel,
	}

	// In production, initialize Raft here
	// For MVP, we use in-memory state

	return e, nil
}

// Shutdown stops the sync engine
func (e *Engine) Shutdown() {
	e.cancel()
	e.logger.Info("Sync engine shutdown complete")
}

// Get retrieves a value by key
func (e *Engine) Get(key string) (interface{}, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if state, exists := e.state[key]; exists {
		return state.Value, nil
	}
	return nil, ErrKeyNotFound
}

// Set stores a value with the given key
func (e *Engine) Set(key string, value interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	version := int64(1)
	if existing, exists := e.state[key]; exists {
		version = existing.Version + 1
	}

	state := &State{
		Key:       key,
		Value:     value,
		Version:   version,
		UpdatedAt: time.Now(),
		UpdatedBy: e.config.NodeID,
	}

	e.state[key] = state

	// Notify watchers
	if watchers, exists := e.watchers[key]; exists {
		for _, w := range watchers {
			go w(key, value, version)
		}
	}

	// In production, replicate via Raft here

	return nil
}

// Delete removes a key
func (e *Engine) Delete(key string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.state[key]; !exists {
		return ErrKeyNotFound
	}

	delete(e.state, key)
	return nil
}

// Watch registers a watcher for a key
func (e *Engine) Watch(key string, watcher Watcher) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.watchers[key] = append(e.watchers[key], watcher)
}

// List returns all keys matching a prefix
func (e *Engine) List(prefix string) ([]*State, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var result []*State
	for key, state := range e.state {
		if len(prefix) == 0 || hasPrefix(key, prefix) {
			result = append(result, state)
		}
	}
	return result, nil
}

// GetVersion returns the current version of a key
func (e *Engine) GetVersion(key string) (int64, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if state, exists := e.state[key]; exists {
		return state.Version, nil
	}
	return 0, ErrKeyNotFound
}

// CompareAndSwap atomically updates a value if the version matches
func (e *Engine) CompareAndSwap(key string, expectedVersion int64, value interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	state, exists := e.state[key]
	if !exists {
		return ErrKeyNotFound
	}

	if state.Version != expectedVersion {
		return ErrVersionMismatch
	}

	state.Value = value
	state.Version++
	state.UpdatedAt = time.Now()
	state.UpdatedBy = e.config.NodeID

	return nil
}

// Snapshot creates a snapshot of the current state
func (e *Engine) Snapshot() ([]byte, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	return json.Marshal(e.state)
}

// Restore restores state from a snapshot
func (e *Engine) Restore(data []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	return json.Unmarshal(data, &e.state)
}

func hasPrefix(s, prefix string) bool {
	if len(prefix) > len(s) {
		return false
	}
	return s[:len(prefix)] == prefix
}

// Errors
var (
	ErrKeyNotFound     = &SyncError{Code: "KEY_NOT_FOUND", Message: "key not found"}
	ErrVersionMismatch = &SyncError{Code: "VERSION_MISMATCH", Message: "version mismatch"}
	ErrNotLeader       = &SyncError{Code: "NOT_LEADER", Message: "not the leader"}
)

// SyncError represents a sync error
type SyncError struct {
	Code    string
	Message string
}

func (e *SyncError) Error() string {
	return e.Message
}
