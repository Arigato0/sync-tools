package cmds

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
	ARGT_ARRAY
)

const (
	_RESET  = "\033[0m"
	RED     = "\033[31m"
	GREEN   = "\033[32m"
	YELLOW  = "\033[33m"
	BLUE    = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN    = "\033[36m"
	GRAY    = "\033[37m"
	WHITE   = "\033[97m"
)

func ColorAs(color string, str string) string {
	return color + str + _RESET
}

type Context struct {
	Args    []any
	Handler *CommandHandler
}

type CommandCallback func(ctx *Context)

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
	AppName      string
	AppNameColor string
	cmdMap       map[string]Command
	aliasMap     map[string]*Command
	ShouldRun    bool
}

func helpCommand(ctx *Context) {
	if len(ctx.Args) > 0 {
		name := ctx.Args[0].(string)
		cmd := ctx.Handler.FindCommand(name)

		if cmd != nil {
			fmt.Println(CommandString(name, cmd))
			return
		}
	}

	ctx.Handler.ShowHelp()

	return
}

func NewCommandHandler() *CommandHandler {
	handler := CommandHandler{}

	handler.AppNameColor = GREEN
	handler.cmdMap = make(map[string]Command)
	handler.aliasMap = make(map[string]*Command)
	handler.ShouldRun = true

	handler.Register("help", Command{
		Alias:       "h",
		Description: "Displays a list of all commands or a single command if given as an argument",
		ArgTypes:    []int{ARGT_STRING},
		Callback:    helpCommand,
	}).Register("quit", Command{
		Alias:       "q",
		Description: "quits the application",
		Callback: func(ctx *Context) {
			ctx.Handler.ShouldRun = false
		},
	})

	return &handler
}

func ArgtToString(argt int, colored bool) string {
	var out string

	switch argt {
	case ARGT_ANY:
		out = "any"
	case ARGT_BOOL:
		out = "true|false"
	case ARGT_STRING:
		out = "string"
	case ARGT_INT:
		out = "integer"
	case ARGT_ARRAY:
		out = "string..."
	default:
		out = "Unknown"
	}

	if colored {
		out = MAGENTA + out + _RESET
	}

	return out
}

func GetFormattedArgt(argTypes []int, minimum int) string {

	builder := strings.Builder{}

	startByte := '<'
	endByte := '>'

	for i, argt := range argTypes {
		toStr := ArgtToString(argt, true)

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
		aliasStr = ColorAs(BLUE, cmd.Alias) + " | "
	}

	return fmt.Sprintf("%s%s: %s\n\tUsage: %s %s",
		aliasStr, ColorAs(BLUE, name), cmd.Description, name, GetFormattedArgt(cmd.ArgTypes, cmd.MinimumArgs))
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

func parseAny(value any) any {
	strValue := value.(string)

	if result, err := strconv.Atoi(strValue); err == nil {
		return result
	} else if result, err := strconv.ParseBool(strValue); err == nil {
		return result
	} else {
		return strValue
	}
}

func createArgs(args []string, cmd *Command) ([]any, error) {
	out := make([]any, len(args))
	argtLen := len(cmd.ArgTypes)

	if argtLen > 0 && cmd.ArgTypes[0] == ARGT_ARRAY {
		for i, arg := range args {
			out[i] = arg
		}

		return out, nil
	}

	for i, arg := range args {

		var value any
		var err error

		argt := -1

		if argtLen > 0 && i < argtLen {
			argt = cmd.ArgTypes[i]
		}

		switch argt {
		case ARGT_INT:
			value, err = strconv.Atoi(arg)
		case ARGT_BOOL:
			value, err = strconv.ParseBool(arg)
		case ARGT_ANY:
			value = parseAny(arg)
		case ARGT_STRING:
			fallthrough
		default:
			value = arg
		}

		if err != nil {
			return nil, err
		}

		out[i] = value
	}

	return out, nil
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

	finalArgs, err := createArgs(args, cmd)

	if err != nil {
		return err
	}

	ctx := Context{
		Args:    finalArgs,
		Handler: handler,
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

func ParseArgs(str string) ([]string, error) {
	parts := strings.Split(str, " ")
	args := make([]string, 0, len(parts))

	for i := 0; i < len(parts); i++ {
		part := parts[i]

		if strings.HasPrefix(part, "\"") || strings.HasPrefix(part, "'") {

			startByte := part[0]

			builder := strings.Builder{}

			part = part[1:]

			for !strings.HasSuffix(part, string(startByte)) {

				builder.WriteString(part + " ")

				i++

				if i >= len(parts) {
					return nil, errors.New("unclosed string")
				}

				part = parts[i]
			}

			part = part[:len(part)-1]
			part = strings.TrimRight(part, " ")

			builder.WriteString(part)

			part = builder.String()
		}

		args = append(args, part)
	}

	return args, nil
}

func (handler *CommandHandler) ExecFromStdin() bool {

	fmt.Printf("[%s]$ ", ColorAs(handler.AppNameColor, handler.AppName))

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

	args := []string{name}

	if rawArgs != "" {
		parsed, err := ParseArgs(rawArgs)

		if err != nil {
			fmt.Println(err.Error())
			return true
		}

		args = append(args, parsed...)
	}

	err = handler.Exec(args)

	if err != nil {
		fmt.Println(err.Error())
	}

	return true
}
