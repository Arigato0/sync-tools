package main

import (
	"fmt"
	"os"

	"lib/cmds"

	"golang.design/x/clipboard"
)

func runFromSysArgs(handler *cmds.CommandHandler) {
	args := os.Args[1:]

	err := handler.Exec(args)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {

	cmdHandler := cmds.NewCommandHandler()

	cmdHandler.AppName = "nclip"

	cmdHandler.Register("test", cmds.Command{
		Description: "A test command",
		MinimumArgs: 2,
		ArgTypes:    []int{cmds.ARGT_STRING, cmds.ARGT_INT},
		Callback: func(ctx *cmds.Context) {
			fmt.Printf("testing with args %s and %d\n", ctx.Args[0], ctx.Args[1])
		},
	}).Register("add", cmds.Command{
		Alias:       "a",
		Description: "Adds the arguments or the top system clipboard entry to the clients nclip database",
		Callback: func(ctx *cmds.Context) {

			if len(ctx.Args) == 0 {
				data := clipboard.Read(clipboard.FmtText)

				if len(data) == 0 {
					fmt.Println(cmds.ColorAs(cmds.RED, "nothing added from clipboard"))
					return
				}

				fmt.Println(string(data))

				fmt.Println("added top clipboard entry to the nclip database")

				return
			}

			for _, arg := range ctx.Args {
				fmt.Println(arg.(string))
			}
		},
	})

	if len(os.Args[1:]) > 0 {
		runFromSysArgs(cmdHandler)
		return
	}

	fmt.Println("use the " + cmds.ColorAs(cmds.YELLOW, "help") + " command for a list of all commands")

	for cmdHandler.ExecFromStdin() {
	}
}
