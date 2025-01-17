package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support/color"

	"github.com/kkumar-gcc/todo/constants"
	"github.com/kkumar-gcc/todo/services"
)

type DeleteTaskCommand struct {
	TaskService services.TaskService
}

// Signature The name and signature of the console command.
func (r *DeleteTaskCommand) Signature() string {
	return "task:delete"
}

// Description The console command description.
func (r *DeleteTaskCommand) Description() string {
	return "Delete an existing task"
}

// Extend The console command extend.
func (r *DeleteTaskCommand) Extend() command.Extend {
	return command.Extend{
		Category: "tasks",
		Flags: []command.Flag{
			&command.IntSliceFlag{
				Name:    "ids",
				Aliases: []string{"i"},
				Usage:   "Comma-separated IDs of the tasks to delete",
			},
		},
	}
}

// Handle Execute the console command.
func (r *DeleteTaskCommand) Handle(ctx console.Context) (err error) {
	taskIDs := ctx.OptionIntSlice("ids")

	if len(taskIDs) == 0 {
		tasks, err := r.TaskService.GetAllTasks(context.Background(), 0, 0, "")
		if err != nil {
			ctx.Error(err.Error())
			return nil
		}

		var choices []console.Choice
		for _, t := range tasks {
			priorityColor := constants.PriorityColors[t.Priority]
			statusColor := constants.StatusColors[t.Status]

			choice := console.Choice{
				Key: color.Sprintf(
					"%s <fg=white;op=bold>(%d) - %s | %s</>",
					t.Title,
					t.ID,
					statusColor,
					priorityColor,
				),
				Value: strconv.Itoa(t.ID),
			}
			choices = append(choices, choice)
		}

		taskStringIDs, err := ctx.MultiSelect("Select the IDs of the tasks to delete:", choices, console.MultiSelectOption{
			Description: "Select the IDs of the tasks to delete",
			Filterable:  true,
			Validate: func(values []string) error {
				if len(values) == 0 {
					return errors.New("at least one task ID is required")
				}
				return nil
			},
		})
		if err != nil {
			ctx.Error(err.Error())
			return nil
		}

		for _, id := range taskStringIDs {
			idInt, err := strconv.Atoi(id)
			if err == nil {
				taskIDs = append(taskIDs, idInt)
			}
		}
	}

	if err := r.TaskService.DeleteTasks(context.Background(), taskIDs); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success(fmt.Sprintf("Task deletion process completed. Removed task IDs: %v", taskIDs))
	return nil
}
