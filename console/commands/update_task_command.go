package commands

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"

	"github.com/kkumar-gcc/todo/constants"
	"github.com/kkumar-gcc/todo/models"
	"github.com/kkumar-gcc/todo/services"
)

type UpdateTaskCommand struct {
	TaskService services.TaskService
}

// Signature The name and signature of the console command.
func (r *UpdateTaskCommand) Signature() string {
	return "task:update"
}

// Description The console command description.
func (r *UpdateTaskCommand) Description() string {
	return "Update an existing task in the TODO application."
}

// Extend The console command extend.
func (r *UpdateTaskCommand) Extend() command.Extend {
	return command.Extend{
		Category: "tasks",
		Flags: []command.Flag{
			&command.IntFlag{
				Name:    "id",
				Aliases: []string{"i"},
				Usage:   "The ID of the task to update",
			},
		},
	}
}

// Handle Execute the console command.
func (r *UpdateTaskCommand) Handle(ctx console.Context) (err error) {
	id := ctx.OptionInt("id")
	if id == 0 {
		idStr, err := ctx.Ask("Enter the ID of the task to update:", console.AskOption{
			Placeholder: "E.g., 1",
			Prompt:      "> ",
			Validate: func(value string) error {
				if strings.TrimSpace(value) == "" {
					return errors.New("task ID is required")
				}
				convInt, convErr := strconv.Atoi(value)
				if convErr != nil {
					return errors.New("task ID must be a valid number")
				}

				if convInt <= 1 {
					return errors.New("task ID must be greater than 1")
				}

				return nil
			},
		})
		if err != nil {
			ctx.Error(err.Error())
			return nil
		}

		id, err = strconv.Atoi(idStr)
		if err != nil {
			ctx.Error("Invalid task ID: " + idStr)
			return nil
		}
	}

	task, err := r.TaskService.GetTaskByID(context.Background(), id)
	if err != nil {
		ctx.Error("Task not found: " + err.Error())
		return nil
	}

	title, err := ctx.Ask("Enter new title for the task:", console.AskOption{
		Placeholder: "E.g., Write article",
		Prompt:      "> ",
		Default:     task.Title,
	})
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	priorityChoices := []console.Choice{
		{Key: constants.PriorityColors[constants.PriorityLow], Value: strconv.Itoa(constants.PriorityLow)},
		{Key: constants.PriorityColors[constants.PriorityMedium], Value: strconv.Itoa(constants.PriorityMedium)},
		{Key: constants.PriorityColors[constants.PriorityHigh], Value: strconv.Itoa(constants.PriorityHigh)},
	}
	priority, err := ctx.Choice("Select priority for the task:", priorityChoices, console.ChoiceOption{
		Default:     strconv.Itoa(task.Priority),
		Description: "Choose a priority for the task",
	})
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	statusChoices := []console.Choice{
		{Key: constants.StatusColors[constants.StatusPending], Value: strconv.Itoa(constants.StatusPending)},
		{Key: constants.StatusColors[constants.StatusInProgress], Value: strconv.Itoa(constants.StatusInProgress)},
		{Key: constants.StatusColors[constants.StatusCompleted], Value: strconv.Itoa(constants.StatusCompleted)},
	}
	status, err := ctx.Choice("Select status for the task:", statusChoices, console.ChoiceOption{
		Default:     strconv.Itoa(task.Status),
		Description: "Choose a status for the task",
	})
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	tags, err := ctx.Ask("Enter tags for the task (comma-separated):", console.AskOption{
		Placeholder: "E.g., work,urgent",
		Prompt:      "> ",
		Default:     task.Tags,
	})
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	priorityInt, err := strconv.Atoi(priority)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	statusInt, err := strconv.Atoi(status)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	err = r.TaskService.UpdateTask(context.Background(), id, func(t *models.Task) (*models.Task, error) {
		t.Title = title
		t.Priority = priorityInt
		t.Status = statusInt
		t.Tags = tags
		return t, nil
	})
	if err != nil {
		ctx.Error("Failed to update task: " + err.Error())
		return nil
	}

	ctx.Success("Task updated successfully!")
	return nil
}
