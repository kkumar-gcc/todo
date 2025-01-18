package main

import (
	"os"

	goravelconsole "github.com/goravel/framework/console"
	"github.com/goravel/framework/support/color"

	"github.com/kkumar-gcc/todo/console"
	"github.com/kkumar-gcc/todo/constants"
)

func main() {
	name := "todo"
	usage := "TODO"
	usageText := "todo [global options] command [command options] [arguments...]"

	cli := goravelconsole.NewApplication(name, usage, usageText, constants.Version, false)

	kernel := &console.Kernel{}

	cli.Register(kernel.Commands())
	if err := cli.Run(os.Args, false); err != nil {
		color.Red().Println(err)
	}
}
