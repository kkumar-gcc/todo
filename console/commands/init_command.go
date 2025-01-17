package commands

import (
	"github.com/kkumar-gcc/todo/database"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

// InitCommand defines the command for initializing the application
type InitCommand struct {
}

// Signature The name and signature of the console command.
func (r *InitCommand) Signature() string {
	return "init"
}

// Description The console command description.
func (r *InitCommand) Description() string {
	return "Set up the necessary components and configurations for the TODO application."
}

// Extend The console command extend.
func (r *InitCommand) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (r *InitCommand) Handle(ctx console.Context) (err error) {
	err = database.RunMigration()
	if err != nil {
		ctx.Error(err.Error())
		return err
	}

	ctx.Success("Application setup completed successfully.")
	return nil
}
