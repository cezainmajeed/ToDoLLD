package repository

import (
	"sync"

	"github.com/cezain/todo-task-tracker/model"
)

type ActivityRepository interface {
	Add(event model.ActivityEvent)
	List() []model.ActivityEvent
}

type InMemmoryActivityRepository struct {
	mu     sync.RWMutex
	events []model.ActivityEvent
}

func NewInMemmoryActivityRepository() *InMemmoryActivityRepository {
	return &InMemmoryActivityRepository{
		events: make([]model.ActivityEvent, 0),
	}
}

func (r *InMemmoryActivityRepository) Add(event model.ActivityEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.events = append(r.events, event)
}

func (r *InMemmoryActivityRepository) List() []model.ActivityEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.ActivityEvent, len(r.events))

	copy(result, r.events)

	return result
}
