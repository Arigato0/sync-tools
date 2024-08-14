package main

import (
	"fmt"
	"lib/cmds"
	"os"
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
	})

	if len(os.Args[1:]) > 0 {
		runFromSysArgs(cmdHandler)
		return
	}

	fmt.Println("use the " + cmds.ColorAs(cmds.YELLOW, "help") + " command for a list of all commands")

	for cmdHandler.ExecFromStdin() {
	}
}
