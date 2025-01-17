package commands

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"

	"github.com/kkumar-gcc/todo/constants"
	"github.com/kkumar-gcc/todo/services"
)

type AddTaskCommand struct {
	TaskService services.TaskService
}

// Signature The name and signature of the console command.
func (r *AddTaskCommand) Signature() string {
	return "task:add"
}

// Description The console command description.
func (r *AddTaskCommand) Description() string {
	return "Create a new task"
}

// Extend The console command extend.
func (r *AddTaskCommand) Extend() command.Extend {
	return command.Extend{
		Category: "tasks",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "title",
				Aliases: []string{"t"},
				Usage:   "The title of the task",
			},
			&command.StringFlag{
				Name:    "priority",
				Aliases: []string{"p"},
				Usage:   "The priority of the task (low, medium, high)",
			},
			&command.StringFlag{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "The status of the task (pending, in-progress, completed)",
			},
			&command.StringFlag{
				Name:    "tags",
				Aliases: []string{"g"},
				Usage:   "Tags for the task, separated by commas",
			},
		},
	}
}

// Handle Execute the console command.
func (r *AddTaskCommand) Handle(ctx console.Context) (err error) {
	title := ctx.Option("title")
	priority := ctx.Option("priority")
	status := ctx.Option("status")
	tags := ctx.Option("tags")

	if title == "" {
		title, err = ctx.Ask("What is the title of the task?", console.AskOption{
			Placeholder: "E.g., Write article",
			Prompt:      "> ",
			Validate: func(value string) error {
				if strings.TrimSpace(value) == "" {
					return errors.New("the task title is required")
				}
				return nil
			},
		})
		if err != nil {
			ctx.Error(err.Error())
			return nil
		}
	}

	if priority == "" {
		choices := []console.Choice{
			{Key: constants.PriorityColors[constants.PriorityLow], Value: strconv.Itoa(constants.PriorityLow)},
			{Key: constants.PriorityColors[constants.PriorityMedium], Value: strconv.Itoa(constants.PriorityMedium)},
			{Key: constants.PriorityColors[constants.PriorityHigh], Value: strconv.Itoa(constants.PriorityHigh)},
		}
		priority, err = ctx.Choice("Select the priority of the task:", choices, console.ChoiceOption{
			Default:     strconv.Itoa(constants.PriorityLow),
			Description: "Choose a priority for the task",
		})
		if err != nil {
			ctx.Error(err.Error())
			return nil
		}
	}

	if status == "" {
		choices := []console.Choice{
			{Key: constants.StatusColors[constants.StatusPending], Value: strconv.Itoa(constants.StatusPending)},
			{Key: constants.StatusColors[constants.StatusInProgress], Value: strconv.Itoa(constants.StatusInProgress)},
			{Key: constants.StatusColors[constants.StatusCompleted], Value: strconv.Itoa(constants.StatusCompleted)},
		}
		status, err = ctx.Choice("Select the status of the task:", choices, console.ChoiceOption{
			Default:     strconv.Itoa(constants.StatusPending),
			Description: "Choose a priority for the task",
		})
		if err != nil {
			ctx.Error(err.Error())
			return nil
		}
	}

	if tags == "" {
		tags, err = ctx.Ask("Enter tags for the task (comma-separated):", console.AskOption{
			Placeholder: "E.g., work,urgent",
			Prompt:      "> ",
		})
		if err != nil {
			ctx.Error(err.Error())
			return nil
		}
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

	if err := r.TaskService.CreateTask(context.Background(), title, statusInt, priorityInt, tags); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Task created successfully!")
	return nil
}
