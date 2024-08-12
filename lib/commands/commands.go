package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ARGT_ANY = iota
	ARGT_INT
	ARGT_STRING
	ARGT_BOOL
)

type Context struct {
	Args []any
}

type CommandCallback func(ctx *Context) int

// nclip add hello
// nclip add_file file_path
// nclip add_dir some_dir

type Command struct {
	Alias       string
	Description string
	MinimumArgs int
	ArgTypes    []int
	Callback    CommandCallback
}

type CommandHandler struct {
	cmdMap   map[string]Command
	aliasMap map[string]*Command
}

func NewCommandHandler() *CommandHandler {
	handler := CommandHandler{}

	handler.cmdMap = make(map[string]Command)
	handler.aliasMap = make(map[string]*Command)

	handler.Register("help", Command{
		Alias:       "h",
		Description: "Displays a list of all commands or a single command if given as an argument",
		ArgTypes:    []int{ARGT_STRING},
		Callback: func(ctx *Context) int {
			if len(ctx.Args) > 0 {
				name := ctx.Args[0].(string)
				cmd, exists := handler.cmdMap[name]

				if exists {
					fmt.Println(CommandString(name, &cmd))
					return 1
				}
			}

			handler.ShowHelp()

			return 1
		},
	})

	return &handler
}

func ArgtToString(argt int) string {
	switch argt {
	case ARGT_ANY:
		return "any"
	case ARGT_BOOL:
		return "true|false"
	case ARGT_STRING:
		return "string"
	case ARGT_INT:
		return "integer"
	default:
		return "Unknown"
	}
}

func GetFormattedArgt(argTypes []int, minimum int) string {

	builder := strings.Builder{}

	startByte := '<'
	endByte := '>'

	for i, argt := range argTypes {
		toStr := ArgtToString(argt)

		if i >= minimum {
			startByte = '|'
			endByte = '|'
		}

		builder.WriteByte(byte(startByte))
		builder.WriteString(toStr)
		builder.WriteByte(byte(endByte))
		builder.WriteByte(' ')
	}

	return builder.String()
}

func CommandString(name string, cmd *Command) string {

	aliasStr := ""

	if cmd.Alias != "" {
		aliasStr = cmd.Alias + " | "
	}
	return fmt.Sprintf("%s%s: %s\n\tUsage: %s %s",
		aliasStr, name, cmd.Description, name, GetFormattedArgt(cmd.ArgTypes, cmd.MinimumArgs))
}

func (handler *CommandHandler) ShowHelp() {
	fmt.Println(`Arguments surrounded by '<>' are required but arguments surrounded by '||' are optional
Commands: `)

	for name, cmd := range handler.cmdMap {
		str := CommandString(name, &cmd)
		fmt.Println(str)
	}
}

func (handler *CommandHandler) FindCommand(name string) *Command {
	cmd, exists := handler.cmdMap[name]

	// attempt to find the alias
	if !exists {
		cmdPtr, exists := handler.aliasMap[name]

		if !exists {
			return nil
		}

		return cmdPtr
	}

	return &cmd
}

func (handler *CommandHandler) Exec(args []string) error {

	if len(args) == 0 {
		return errors.New("expected args to be greater than 0")
	}

	cmd := handler.FindCommand(args[0])

	if cmd == nil {
		return fmt.Errorf("command '%s' does not exist", args[0])
	}

	args = args[1:]

	argLen := len(args)
	argtLen := len(cmd.ArgTypes)

	if argLen < int(cmd.MinimumArgs) && argLen != argtLen {
		return fmt.Errorf("expected %d args but got %d", cmd.MinimumArgs, argLen)
	}

	ctx := Context{}

	ctx.Args = make([]any, argLen)

	for i, argt := range cmd.ArgTypes {

		if i >= argLen {
			break
		}

		var value any
		var err error

		switch argt {
		case ARGT_STRING:
			fallthrough
		case ARGT_ANY:
			value = args[i]
		case ARGT_INT:
			value, err = strconv.Atoi(args[i])
		case ARGT_BOOL:
			value, err = strconv.ParseBool(args[i])
		}

		if err != nil {
			return err
		}

		ctx.Args[i] = value
	}

	cmd.Callback(&ctx)

	return nil
}

func (cmdHandler *CommandHandler) Register(name string, command Command) *CommandHandler {

	cmdHandler.cmdMap[name] = command

	if command.Alias != "" {
		cmdHandler.aliasMap[command.Alias] = &command
	}

	return cmdHandler
}

func (handler *CommandHandler) ExecFromStdin() bool {
	reader := bufio.NewReader(os.Stdin)

	name, err := reader.ReadString('\n')

	if err != nil {
		return false
	}

	name = strings.Trim(name, "\n ")

	spaceIndex := strings.Index(name, " ")

	rawArgs := ""

	if spaceIndex > -1 {
		rawArgs = name[spaceIndex+1:]
		name = name[:spaceIndex]
	}

	fmt.Println(rawArgs)

	err = handler.Exec([]string{name})

	if err != nil {
		fmt.Println(err.Error())
	}

	return true
}
