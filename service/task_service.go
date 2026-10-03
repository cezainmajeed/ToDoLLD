package service

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cezain/todo-task-tracker/filter"
	"github.com/cezain/todo-task-tracker/model"
	"github.com/cezain/todo-task-tracker/repository"
)

type DefaultTaskService struct {
	taskRepository     repository.TaskRepository
	activityRepository repository.ActivityRepository

	mu        sync.Mutex
	idCounter int
}

func NewTaskService(taskRepository repository.TaskRepository, activityRepository repository.ActivityRepository) TaskService {
	return &DefaultTaskService{
		taskRepository:     taskRepository,
		activityRepository: activityRepository,
	}
}

type TaskService interface {
	AddTask(task *model.Task) error
	GetTask(taskId string, userId string) (*model.Task, error)
	ModifyTask(task *model.Task) error
	RemoveTask(taskID string, userID string) error
	CompleteTask(taskID string, userID string) error
	ListTask(req ListTaskRequest) (*[]model.Task, error)
	GetActivityLog(userID string, period *TimePeriod) ([]model.ActivityEvent, error)
}

func (d *DefaultTaskService) generateId(prefixString string) string {
	d.idCounter++
	return fmt.Sprintf("%s-%d", prefixString, d.idCounter)
}

// AddTask implements TaskService.
func (d *DefaultTaskService) AddTask(task *model.Task) error {
	// generateId
	// check if task already exists
	// create and save task
	// record activity

	if task == nil {
		return errors.New("task cant be null")
	}

	if strings.TrimSpace(task.UserId) == "" {
		return errors.New("userId is requirec")
	}

	if strings.TrimSpace(task.Title) == "" {
		return errors.New("title is requirec")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()

	if task.ID == "" {
		task.ID = d.generateId("TASK")
	}

	if _, err := d.taskRepository.Get(task.ID); err == nil {
		return errors.New("task already exists")
	}

	task.CreatedAt = now
	task.UpdatedAt = now
	task.Status = model.StatusActive

	if err := d.taskRepository.Save(task); err != nil {
		return err
	}

	d.recordActivity(task, model.ActivityAdded, now, fmt.Sprintf("Task %s added", task.ID))

	return nil
}

type TimePeriod struct {
	Start time.Time
	End   time.Time
}

// GetActivityLog implements TaskService.
func (d *DefaultTaskService) GetActivityLog(userID string, period *TimePeriod) ([]model.ActivityEvent, error) {
	// get all activity
	// range and filter
	// return

	events := d.activityRepository.List()

	result := make([]model.ActivityEvent, 0)

	for _, event := range events {
		if event.UserID != userID {
			continue
		}

		if period != nil && (event.Timestamp.Before(period.Start) || event.Timestamp.After(period.End)) {
			continue
		}

		result = append(result, event)
	}

	return result, nil
}

// GetTask implements TaskService.
func (d *DefaultTaskService) GetTask(taskId string, userId string) (*model.Task, error) {
	task, err := d.taskRepository.Get(taskId)
	if err != nil {
		return nil, err
	}

	if task.UserId != userId {
		return nil, errors.New("unauthorised to access")
	}

	return task, nil
}

type ListTaskRequest struct {
	UserID  string
	Filters []filter.TaskFilter
}

// ListTask implements TaskService.
func (d *DefaultTaskService) ListTask(req ListTaskRequest) (*[]model.Task, error) {
	// get all task
	// make empty list
	// check userid and status
	//match for filter
	// return matched

	allTasks := d.taskRepository.List()

	result := make([]model.Task, 0)

	for _, task := range allTasks {

		// not including completed/rejected task in active TODO list
		if task.Status != model.StatusActive {
			continue
		}

		if task.UserId != req.UserID {
			continue
		}

		matches := true

		for _, filter := range req.Filters {
			if filter == nil {
				continue
			}

			if !filter.Match(task) {
				matches = false
				break
			}
		}

		if matches {
			result = append(result, *task)
		}
	}

	return &result, nil
}

// ModifyTask implements TaskService.
func (d *DefaultTaskService) ModifyTask(task *model.Task) error {
	// get task by Id
	// check userId
	// check status
	// update task obj
	// save task
	// recored activity

	if task == nil {
		return errors.New("task cant be null")
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	existingTask, err := d.taskRepository.Get(task.ID)
	if err != nil {
		return err
	}

	if existingTask.UserId != task.UserId {
		return errors.New("Unauthorised access")
	}

	if existingTask.Status != model.StatusActive {
		return errors.New("inactive task cant be modified")
	}

	now := time.Now()

	task.CreatedAt = existingTask.CreatedAt
	task.UpdatedAt = now
	task.Status = existingTask.Status

	if err := d.taskRepository.Save(task); err != nil {
		return err
	}

	d.recordActivity(
		task,
		model.ActivityModified,
		now,
		fmt.Sprintf("Task  %s modified", task.ID),
	)

	return nil
}

// RemoveTask implements TaskService.
func (d *DefaultTaskService) RemoveTask(taskID string, userID string) error {
	//get task
	//check userId
	//check status
	//update status and save
	//record activity

	d.mu.Lock()
	defer d.mu.Unlock()

	task, err := d.taskRepository.Get(taskID)
	if err != nil {
		return err
	}

	if task.UserId != userID {
		return errors.New("Unauthorised access")
	}

	if task.Status != model.StatusActive {
		return errors.New("inactive task cant be modified")
	}

	now := time.Now()

	task.Status = model.StatusRemoved
	task.RemovedAt = &now
	task.UpdatedAt = now

	if err := d.taskRepository.Save(task); err != nil {
		return err
	}

	d.recordActivity(
		task,
		model.ActivityRemoved,
		now,
		fmt.Sprintf("Task  %s removed", task.ID),
	)

	return nil
}

// CompleteTask implements TaskService.
func (d *DefaultTaskService) CompleteTask(taskID string, userID string) error {
	//get task
	//check userId
	//check status
	//update status and save
	//record activity

	d.mu.Lock()
	defer d.mu.Unlock()

	task, err := d.taskRepository.Get(taskID)
	if err != nil {
		return err
	}

	if task.UserId != userID {
		return errors.New("Unauthorised access")
	}

	if task.Status != model.StatusActive {
		return errors.New("inactive task cant be modified")
	}

	now := time.Now()

	task.Status = model.StatusCompleted
	task.RemovedAt = &now
	task.UpdatedAt = now

	if err := d.taskRepository.Save(task); err != nil {
		return err
	}

	d.recordActivity(
		task,
		model.ActivityCompleted,
		now,
		fmt.Sprintf("Task  %s completed", task.ID),
	)

	return nil
}

func (d *DefaultTaskService) recordActivity(
	task *model.Task,
	activityType model.ActivityType,
	timeStamp time.Time,
	details string,
) {
	event := model.ActivityEvent{
		ID:        d.generateId("EVENT"),
		TaskID:    task.ID,
		UserID:    task.UserId,
		Type:      activityType,
		Timestamp: timeStamp,
		Details:   details,
	}

	d.activityRepository.Add(event)
}
