package console

import (
	"github.com/goravel/framework/contracts/console"

	"github.com/kkumar-gcc/todo/console/commands"
	"github.com/kkumar-gcc/todo/database"
	"github.com/kkumar-gcc/todo/repositories"
	"github.com/kkumar-gcc/todo/services"
)

type Kernel struct {
}

func (kernel *Kernel) Commands() []console.Command {
	taskRepository := repositories.NewTaskRepository(database.GetInstance())
	taskService := services.NewTaskService(taskRepository)
	return []console.Command{
		&commands.AddTaskCommand{
			TaskService: taskService,
		},
		&commands.ListTasksCommand{
			TaskService: taskService,
		},
		&commands.DeleteTaskCommand{
			TaskService: taskService,
		},
		&commands.UpdateTaskCommand{
			TaskService: taskService,
		},
		&commands.InitCommand{},
	}
}
