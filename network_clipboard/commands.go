package main

import (
	"fmt"
	"lib/cmds"
	clipdb "network_clipboard/clip_db"

	"golang.design/x/clipboard"
)

func addFs(entryType int, args []any) {
	path := args[0].(string)

	entry := clipdb.NewFsEntry(path, entryType)
	err := entry.Save([]byte(path))

	if err != nil {
		fmt.Println(cmds.ColorAs(cmds.RED, err.Error()))
	} else {
		fmt.Println("copied directory to nclip database")
	}
}

func AddDirCommand(ctx *cmds.Context) {
	addFs(clipdb.TYPE_DIR, ctx.Args)
}

func AddFileCommand(ctx *cmds.Context) {
	addFs(clipdb.TYPE_FILE, ctx.Args)
}

func AddCommand(ctx *cmds.Context) {

	if len(ctx.Args) == 0 {
		data := clipboard.Read(clipboard.FmtText)

		if len(data) == 0 {
			fmt.Println(cmds.ColorAs(cmds.RED, "nothing added from clipboard"))
			return
		}

		entry := clipdb.NewTextEntry()
		err := entry.Save(data)

		if err != nil {
			fmt.Println(cmds.ColorAs(cmds.RED, err.Error()))
		} else {
			fmt.Println("added top clipboard item to the nclip database")
		}

		return
	}

	for _, arg := range ctx.Args {
		fmt.Println(arg.(string))
	}
}
