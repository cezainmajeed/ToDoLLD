package model

import "time"

type TaskStatus string

const (
	StatusActive    TaskStatus = "ACTIVE"
	StatusCompleted TaskStatus = "COMPLETED"
	StatusRejected  TaskStatus = "REJECTED"
	StatusRemoved   TaskStatus = "REMOVED"
)

type Task struct {
	ID     string
	UserId string

	Title       string
	Description string

	Tags map[string]struct{}

	Status TaskStatus

	Deadline *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time

	CompletedAt *time.Time
	RemovedAt   *time.Time
}

func (t *Task) Clone() *Task {
	if t == nil {
		return nil
	}

	copyTask := *t

	if t.Deadline != nil {
		deadLine := *t.Deadline
		copyTask.Deadline = &deadLine
	}

	if t.CompletedAt != nil {
		completedAt := *t.CompletedAt
		copyTask.CompletedAt = &completedAt
	}

	if t.RemovedAt != nil {
		removedAt := *t.RemovedAt
		copyTask.RemovedAt = &removedAt
	}

	copyTask.Tags = make(map[string]struct{})

	for tag := range t.Tags {
		copyTask.Tags[tag] = struct{}{}
	}

	return &copyTask
}
