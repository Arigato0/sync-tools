package main

import (
	"fmt"
	"os"

	"lib/commands"
)

func runFromSysArgs(handler *commands.CommandHandler) {
	args := os.Args[1:]

	err := handler.Exec(args)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func main() {

	cmdHandler := commands.NewCommandHandler()

	cmdHandler.Register("test", commands.Command{
		Description: "A test command",
		MinimumArgs: 2,
		ArgTypes:    []int{commands.ARGT_STRING, commands.ARGT_INT},
		Callback: func(ctx *commands.Context) int {
			fmt.Printf("testing with args %s and %d\n", ctx.Args[0], ctx.Args[1])
			return 0
		},
	})

	if len(os.Args[1:]) > 0 {
		runFromSysArgs(cmdHandler)
		return
	}

	shouldRun := true

	for shouldRun {
		shouldRun = cmdHandler.ExecFromStdin()
	}
}
