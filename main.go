package main

import (
	"fmt"
	"time"

	"github.com/cezain/todo-task-tracker/filter"
	"github.com/cezain/todo-task-tracker/model"
	"github.com/cezain/todo-task-tracker/repository"
	"github.com/cezain/todo-task-tracker/service"
)

func main() {
	taskRepository := repository.NewInMemmoryTaskRepository()
	activityRepository := repository.NewInMemmoryActivityRepository()

	taskService := service.NewTaskService(
		taskRepository,
		activityRepository,
	)

	fmt.Println("TODO task tracker")

	//add task
	fmt.Println("Adding tasks...")

	deadline := time.Now().Add(24 * time.Hour)

	task1 := &model.Task{
		UserId:      "user-1",
		Title:       "Task title 1",
		Description: "Task description 1",
		Deadline:    &deadline,
		Tags: map[string]struct{}{
			"tag-1": {},
			"tag-2": {},
		},
	}

	if err := taskService.AddTask(task1); err != nil {
		panic(err)
	}

	task2 := &model.Task{
		UserId:      "user-1",
		Title:       "Task title 2",
		Description: "Task description 2",
		Deadline:    &deadline,
		Tags: map[string]struct{}{
			"tag-3": {},
			"tag-2": {},
		},
	}

	if err := taskService.AddTask(task2); err != nil {
		panic(err)
	}

	task3 := &model.Task{
		UserId: "user-2",
		Title:  "Task title 3",
		Tags: map[string]struct{}{
			"tag-5": {},
		},
	}

	if err := taskService.AddTask(task3); err != nil {
		panic(err)
	}

	//list task
	fmt.Println("listing tasks")
	tasks, err := taskService.ListTask(
		service.ListTaskRequest{
			UserID: "user-1",
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(tasks)

	//filter
	fmt.Println("listing tasks with filter")
	tasks, err = taskService.ListTask(
		service.ListTaskRequest{
			UserID: "user-1",
			Filters: []filter.TaskFilter{
				filter.TagFilter{
					Tag: "tag-2",
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(tasks)

	//modify
	fmt.Println("Modiying tasks")

	task1.Title = "modified title 1"
	task1.Description = "modified desc 1"
	task1.Tags = map[string]struct{}{
		"tag-5": {},
		"tag-6": {},
	}

	if err := taskService.ModifyTask(task1); err != nil {
		panic(err)
	}

	updatedTask, err := taskService.GetTask(task1.ID, "user-1")
	if err != nil {
		panic(err)
	}

	fmt.Println(updatedTask)

	// complete

	fmt.Println("Completing tasks")

	if err := taskService.CompleteTask(task2.ID, "user-1"); err != nil {
		panic(err)
	}

	completedTask, err := taskService.GetTask(task2.ID, "user-1")
	if err != nil {
		panic(err)
	}

	fmt.Println(completedTask)

	//remove

	fmt.Println("Removing tasks")

	if err := taskService.RemoveTask(task1.ID, "user-1"); err != nil {
		panic(err)
	}

	removedTask, err := taskService.GetTask(task1.ID, "user-1")
	if err != nil {
		panic(err)
	}

	fmt.Println(removedTask)

	// activity logs

	fmt.Println("Getting activity logs")

	events, err := taskService.GetActivityLog("user-1", nil)

	if err != nil {
		panic(err)
	}

	fmt.Println(events)

}
