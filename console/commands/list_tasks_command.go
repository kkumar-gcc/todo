package commands

import (
	"context"
	"time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support/color"

	"github.com/kkumar-gcc/todo/constants"
	"github.com/kkumar-gcc/todo/models"
	"github.com/kkumar-gcc/todo/services"
)

type ListTasksCommand struct {
	TaskService services.TaskService
}

// Signature The name and signature of the console command.
func (r *ListTasksCommand) Signature() string {
	return "task:list"
}

// Description The console command description.
func (r *ListTasksCommand) Description() string {
	return "List all tasks"
}

// Extend The console command extend.
func (r *ListTasksCommand) Extend() command.Extend {
	return command.Extend{
		Category: "tasks",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "sort",
				Aliases: []string{"s"},
				Usage:   "Sort tasks by field (status or priority)",
			},
			&command.StringFlag{
				Name:    "status",
				Aliases: []string{"st"},
				Usage:   "Filter tasks by status (pending, in_progress, completed)",
			},
			&command.StringFlag{
				Name:    "priority",
				Aliases: []string{"p"},
				Usage:   "Filter tasks by priority (low, medium, high)",
			},
		},
	}
}

// Handle Execute the console command.
func (r *ListTasksCommand) Handle(ctx console.Context) (err error) {
	sort, status, priority := ctx.Option("sort"), ctx.Option("status"), ctx.Option("priority")
	tasks, err := r.TaskService.GetAllTasks(context.Background(), constants.StatusMap[status], constants.PriorityMap[priority], sort)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if len(tasks) == 0 {
		ctx.Info("No tasks found matching the given criteria.")
		return nil
	}

	ctx.NewLine()
	color.Println("<fg=blue;op=bold>Task List:</>")
	ctx.NewLine()

	groupedTasks := r.groupTasks(tasks, sort)
	for label, tasksGroup := range groupedTasks {
		ctx.TwoColumnDetail(color.Sprintf("<fg=cyan;op=bold>%s</>", label), "Details")
		for _, task := range tasksGroup {
			idLabel := color.Sprintf("<fg=white;op=bold>%d</>", task.ID)
			statusLabel := constants.StatusColors[task.Status]
			priorityLabel := constants.PriorityColors[task.Priority]
			tagsAndCreatedAt := color.Sprintf("<fg=gray>Tags: %s, Created At: %s</>", task.Tags, task.CreatedAt.Format(time.RFC822))
			ctx.TwoColumnDetail(task.Title+" ("+idLabel+") "+tagsAndCreatedAt, statusLabel+" | "+priorityLabel)
		}
		ctx.NewLine()
	}

	return nil
}

func (r *ListTasksCommand) groupTasks(tasks []models.Task, sort string) map[string][]models.Task {
	grouped := make(map[string][]models.Task)
	for _, task := range tasks {
		var key string
		if sort == "priority" {
			key = constants.PriorityLabels[task.Priority]
		} else if sort == "status" {
			key = constants.StatusLabels[task.Status]
		} else {
			key = "Tasks"
		}
		grouped[key] = append(grouped[key], task)
	}
	return grouped
}
