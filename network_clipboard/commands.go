package main

import (
	"encoding/json"
	"fmt"
	"lib/cmds"
	clipdb "network_clipboard/clip_db"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

	var data []byte

	if len(ctx.Args) == 0 {
		data = clipboard.Read(clipboard.FmtText)

		if len(data) == 0 {
			fmt.Println(cmds.ColorAs(cmds.RED, "nothing added from clipboard"))
			return
		}
	} else {
		builder := strings.Builder{}

		for i, arg := range ctx.Args {
			builder.WriteString(arg.(string))

			if i < len(ctx.Args)-1 {
				builder.WriteRune(' ')
			}
		}

		data = []byte(builder.String())
	}

	entry := clipdb.NewTextEntry()
	err := entry.Save(data)

	if err != nil {
		fmt.Println(cmds.ColorAs(cmds.RED, err.Error()))
	} else {
		fmt.Println("added top clipboard item to the nclip database")
	}
}

func truncate(s string, maxLen int) string {
	if len(s) < maxLen {
		return s
	}

	return s[:maxLen] + "..."
}

func ViewCommand(ctx *cmds.Context) {
	dbDir := clipdb.GetNclipDbDir()

	entryFiles, err := os.ReadDir(dbDir)

	if err != nil {
		return
	}

	entries := make([]clipdb.Entry, 0, len(entryFiles)/2)

	for _, entryFile := range entryFiles {
		name := entryFile.Name()

		if !strings.HasSuffix(name, clipdb.NCLIP_ENTRY_EXT) {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dbDir, name))

		if err != nil {
			continue
		}

		var entry clipdb.Entry

		err = json.Unmarshal(data, &entry)

		if err != nil {
			continue
		}

		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Time.Unix() > entries[j].Time.Unix()
	})

	for i, entry := range entries {
		contents := entry.Filename

		if entry.Type == clipdb.TYPE_TEXT {
			contents = truncate(string(entry.Data), 32)
		}

		fmt.Printf("%d - (%s) %s\n", i, clipdb.TypeString(entry.Type), contents)
	}
}
