package main

import (
	"fmt"
	"lib/cmds"
	"network_clipboard/server"
	"os"
)

// TODO: add a config file that the user can modify through commands

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

	cmdHandler.Register("add", cmds.Command{
		Alias:       "a",
		Description: "Adds the arguments or the top system clipboard entry to the clients nclip database",
		ArgTypes:    []int{cmds.ARGT_ARRAY},
		Callback:    AddCommand,
	}).Register("add_dir", cmds.Command{
		Alias:       "ad",
		Description: "Adds the directory given in the arguments to the nclip database",
		MinimumArgs: 1,
		ArgTypes:    []int{cmds.ARGT_STRING},
		Callback:    AddDirCommand,
	}).Register("add_file", cmds.Command{
		Alias:       "af",
		Description: "Adds the file given in the arguments to the nclip database",
		MinimumArgs: 1,
		ArgTypes:    []int{cmds.ARGT_STRING},
		Callback:    AddFileCommand,
	}).Register("view", cmds.Command{
		Alias:       "v",
		Description: "view all nclip entries either on the local machine if no arguments are given or on a clients machine if given as the first argument",
		ArgTypes:    []int{cmds.ARGT_STRING},
		Callback:    ViewCommand,
	}).Register("test", cmds.Command{
		Alias:       "t",
		Description: "a command used for testing",
		MinimumArgs: 2,
		ArgTypes:    []int{cmds.ARGT_STRING, cmds.ARGT_ANY},
		Callback: func(ctx *cmds.Context) {
			fmt.Printf("%s %t\n", ctx.Args[0], ctx.Args[1])
		},
	}).Register("view_servers", cmds.Command{
		Alias:       "vs",
		Description: "Shows a list of all nclip servers in the local network",
		Callback: func(ctx *cmds.Context) {
			names, err := server.GetHostNames()

			if err != nil {
				fmt.Println(cmds.ColorAs(cmds.RED, err.Error()))
				return
			}

			fmt.Println(names)
		},
	})

	if len(os.Args[1:]) > 0 {
		runFromSysArgs(cmdHandler)
		return
	}

	fmt.Println("use the " + cmds.ColorAs(cmds.YELLOW, "help") + " command for a list of all commands")

	for cmdHandler.ShouldRun && cmdHandler.ExecFromStdin() {
	}
}
