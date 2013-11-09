package main

import (
	"fmt"
	"os"
)

type commander interface {
	Parse([]string)
	Run(*Settings) error
}

var (
	commands = make(map[string]commander)
)

func printHelp(errs ...string) {
	content := `newbeego command usage:

    blank      - create blank app with no config settings
    help       - print this help
`

	if len(errs) > 0 {
		fmt.Println(errs[0])
	}
	fmt.Println(content)
	os.Exit(2)
}

func (s *Settings) RunCommand(params ...string) {
	/*cmdName := "newbeego"
	if len(params) > 0 {
		cmdName = params[0]
	}*/

	// if len(os.Args) < 2 || os.Args[1] != cmdName {
	if len(os.Args) < 1 {
		printHelp(fmt.Sprintf("newbeego error: no command specified"))
		return
	}

	// args := argString(os.Args[2:])
	args := argString(os.Args[1:])
	name := args.Get(0)

	if name == "help" {
		printHelp()
	}

	if cmd, ok := commands[name]; ok {
		// cmd.Parse(os.Args[3:])
		cmd.Parse(os.Args[2:])
		cmd.Run(s)
		os.Exit(0)
	} else {
		if name == "" {
			printHelp()
		} else {
			printHelp(fmt.Sprintf("unknown command %s", name))
		}
	}
}

func init() {
	commands["blank"] = new(newBlank)
}
