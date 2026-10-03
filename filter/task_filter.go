package filter

import (
	"strings"
	"time"

	"github.com/cezain/todo-task-tracker/model"
)

type TaskFilter interface {
	Match(task *model.Task) bool
}

// filter by tag, status, deadline

// filter by tag
type TagFilter struct {
	Tag string
}

func (f TagFilter) Match(task *model.Task) bool {
	tag := strings.ToLower(strings.TrimSpace(f.Tag))

	_, exists := task.Tags[tag]

	return exists
}

// filter by deadline
type DeadlineBeforeFilter struct {
	Time time.Time
}

func (f DeadlineBeforeFilter) Match(task *model.Task) bool {
	if task.Deadline == nil {
		return false
	}

	return !task.Deadline.After(f.Time)
}

// filter nby status
type StatusFilter struct {
	Status model.TaskStatus
}

func (f StatusFilter) Match(task *model.Task) bool {
	return task.Status == f.Status
}
