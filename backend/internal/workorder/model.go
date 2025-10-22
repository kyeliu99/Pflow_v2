package workorder

import "time"

type Status string

type WorkOrder struct {
	ID        string            `json:"id" db:"id"`
	FlowID    string            `json:"flowId" db:"flow_id"`
	Title     string            `json:"title" db:"title"`
	Assignee  string            `json:"assignee" db:"assignee"`
	Status    Status            `json:"status" db:"status"`
	Payload   map[string]any    `json:"payload" db:"payload"`
	Metadata  map[string]string `json:"metadata" db:"metadata"`
	CreatedAt time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time         `json:"updatedAt" db:"updated_at"`
}

const (
	StatusPending  Status = "pending"
	StatusRunning  Status = "running"
	StatusFailed   Status = "failed"
	StatusComplete Status = "complete"
)
