// Package protocol implements the ACIP (Agent Collaboration Interaction Protocol)
package protocol

import (
	"encoding/json"
	"time"
)

// MessageType defines the type of ACIP message
type MessageType string

const (
	// MsgTypeRequest is for task requests and assignments
	MsgTypeRequest MessageType = "request"
	// MsgTypeReview is for code/work review messages
	MsgTypeReview MessageType = "review"
	// MsgTypeSync is for state synchronization messages
	MsgTypeSync MessageType = "sync"
	// MsgTypeFeedback is for feedback and results
	MsgTypeFeedback MessageType = "feedback"
	// MsgTypeApproval is for approval/decision messages
	MsgTypeApproval MessageType = "approval"
)

// Priority defines message priority levels
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityNormal   Priority = "normal"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// Message is the base ACIP message structure
type Message struct {
	// Header contains message metadata
	Header Header `json:"header"`

	// Payload contains the actual message content
	Payload Payload `json:"payload"`

	// Context provides additional context for the message
	Context Context `json:"context,omitempty"`

	// TraceID for distributed tracing
	TraceID string `json:"traceId,omitempty"`
}

// Header contains message metadata
type Header struct {
	// Message ID (unique identifier)
	ID string `json:"id"`

	// Message type
	Type MessageType `json:"type"`

	// Message priority
	Priority Priority `json:"priority"`

	// Sender agent ID
	From string `json:"from"`

	// Recipient agent ID(s), empty for broadcast
	To []string `json:"to,omitempty"`

	// Timestamp
	Timestamp time.Time `json:"timestamp"`

	// Correlation ID for request-response pairing
	CorrelationID string `json:"correlationId,omitempty"`

	// TTL in seconds (0 = no expiry)
	TTL int `json:"ttl,omitempty"`
}

// Payload contains the actual message content
type Payload struct {
	// Action to be performed
	Action string `json:"action"`

	// Content/parameters of the action
	Content json.RawMessage `json:"content"`

	// Attachments (files, URLs, etc.)
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents a file or resource attachment
type Attachment struct {
	Type        string `json:"type"` // file, url, data
	Name        string `json:"name"`
	URL         string `json:"url,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Size        int64  `json:"size,omitempty"`
}

// Context provides additional context for the message
type Context struct {
	// Task ID this message relates to
	TaskID string `json:"taskId,omitempty"`

	// Workflow ID
	WorkflowID string `json:"workflowId,omitempty"`

	// Step/Stage in the workflow
	Step string `json:"step,omitempty"`

	// Parent message ID (for threaded conversations)
	ParentID string `json:"parentId,omitempty"`

	// Custom metadata
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// RequestMessage is a task request message
type RequestMessage struct {
	TaskID      string                 `json:"taskId"`
	TaskType    string                 `json:"taskType"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
	Deadline    *time.Time             `json:"deadline,omitempty"`
	Requirements []Requirement         `json:"requirements,omitempty"`
}

// Requirement specifies a task requirement
type Requirement struct {
	Type     string   `json:"type"`     // skill, resource, permission
	Name     string   `json:"name"`
	Value    string   `json:"value"`
	Optional bool     `json:"optional"`
	Tags     []string `json:"tags,omitempty"`
}

// ReviewMessage is a review/feedback message
type ReviewMessage struct {
	TargetID    string       `json:"targetId"`    // What is being reviewed
	TargetType  string       `json:"targetType"`  // code, design, document
	Status      ReviewStatus `json:"status"`
	Comments    []Comment    `json:"comments,omitempty"`
	Suggestions []string     `json:"suggestions,omitempty"`
}

// ReviewStatus represents the status of a review
type ReviewStatus string

const (
	ReviewStatusApproved   ReviewStatus = "approved"
	ReviewStatusChanges    ReviewStatus = "changes_requested"
	ReviewStatusRejected   ReviewStatus = "rejected"
	ReviewStatusPending    ReviewStatus = "pending"
)

// Comment represents a review comment
type Comment struct {
	ID        string `json:"id"`
	Author    string `json:"author"`
	Content   string `json:"content"`
	Location  string `json:"location,omitempty"` // Line number, file path, etc.
	Timestamp time.Time `json:"timestamp"`
}

// SyncMessage is a state synchronization message
type SyncMessage struct {
	ResourceType string      `json:"resourceType"` // task, agent, workflow
	ResourceID   string      `json:"resourceId"`
	Operation    SyncOp      `json:"operation"`
	Data         interface{} `json:"data"`
	Version      int64       `json:"version"`
}

// SyncOp represents sync operation type
type SyncOp string

const (
	SyncOpCreate SyncOp = "create"
	SyncOpUpdate SyncOp = "update"
	SyncOpDelete SyncOp = "delete"
)

// FeedbackMessage is a feedback/result message
type FeedbackMessage struct {
	TaskID    string                 `json:"taskId"`
	Status    FeedbackStatus         `json:"status"`
	Result    interface{}            `json:"result,omitempty"`
	Error     *ErrorDetail           `json:"error,omitempty"`
	Metrics   *ExecutionMetrics      `json:"metrics,omitempty"`
	Artifacts []Artifact             `json:"artifacts,omitempty"`
}

// FeedbackStatus represents the status of feedback
type FeedbackStatus string

const (
	FeedbackStatusSuccess FeedbackStatus = "success"
	FeedbackStatusFailure FeedbackStatus = "failure"
	FeedbackStatusPartial FeedbackStatus = "partial"
)

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ExecutionMetrics contains execution metrics
type ExecutionMetrics struct {
	Duration    time.Duration `json:"duration"`
	TokenUsage  int64         `json:"tokenUsage"`
	APICalls    int           `json:"apiCalls"`
	RetryCount  int           `json:"retryCount"`
}

// Artifact represents a generated artifact
type Artifact struct {
	Type        string `json:"type"` // file, url, data
	Name        string `json:"name"`
	URL         string `json:"url,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Size        int64  `json:"size,omitempty"`
}

// ApprovalMessage is an approval/decision message
type ApprovalMessage struct {
	RequestID    string        `json:"requestId"`
	RequestType  string        `json:"requestType"` // task, change, access
	Status       ApprovalStatus `json:"status"`
	Approver     string        `json:"approver"`
	Reason       string        `json:"reason,omitempty"`
	Conditions   []string      `json:"conditions,omitempty"`
	ExpiresAt    *time.Time    `json:"expiresAt,omitempty"`
}

// ApprovalStatus represents the status of an approval
type ApprovalStatus string

const (
	ApprovalStatusApproved ApprovalStatus = "approved"
	ApprovalStatusDenied   ApprovalStatus = "denied"
	ApprovalStatusPending  ApprovalStatus = "pending"
	ApprovalStatusExpired  ApprovalStatus = "expired"
)

// NewMessage creates a new ACIP message
func NewMessage(msgType MessageType, from string, to []string) *Message {
	return &Message{
		Header: Header{
			ID:        generateID(),
			Type:      msgType,
			Priority:  PriorityNormal,
			From:      from,
			To:        to,
			Timestamp: time.Now(),
		},
		Payload: Payload{},
	}
}

// WithPriority sets message priority
func (m *Message) WithPriority(p Priority) *Message {
	m.Header.Priority = p
	return m
}

// WithCorrelationID sets correlation ID
func (m *Message) WithCorrelationID(id string) *Message {
	m.Header.CorrelationID = id
	return m
}

// WithContext sets message context
func (m *Message) WithContext(ctx Context) *Message {
	m.Context = ctx
	return m
}

// WithTraceID sets trace ID for distributed tracing
func (m *Message) WithTraceID(id string) *Message {
	m.TraceID = id
	return m
}

// SetContent sets the payload content
func (m *Message) SetContent(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	m.Payload.Content = data
	return nil
}

// GetContent unmarshals the payload content
func (m *Message) GetContent(v interface{}) error {
	return json.Unmarshal(m.Payload.Content, v)
}

// ToJSON serializes the message to JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON deserializes a message from JSON
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().Nanosecond()%len(letters)]
	}
	return string(b)
}
