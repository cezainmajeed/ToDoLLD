package repository

import (
	"errors"
	"sync"

	"github.com/cezain/todo-task-tracker/model"
)

type TaskRepository interface {
	Get(taskID string) (*model.Task, error)
	Save(task *model.Task) error
	List() []*model.Task
}

type InMemmoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*model.Task
}

func NewInMemmoryTaskRepository() *InMemmoryTaskRepository {
	return &InMemmoryTaskRepository{
		tasks: make(map[string]*model.Task),
	}
}

func (r *InMemmoryTaskRepository) Get(taskID string) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exist := r.tasks[taskID]

	if !exist {
		return nil, errors.New("no task found")
	}

	return task, nil
}

func (r *InMemmoryTaskRepository) Save(task *model.Task) error {
	if task == nil {
		return errors.New("task cant be nil")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = task.Clone()

	return nil
}

func (r *InMemmoryTaskRepository) List() []*model.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*model.Task, 0, len(r.tasks))

	for _, task := range r.tasks {
		tasks = append(tasks, task.Clone())
	}

	return tasks
}
