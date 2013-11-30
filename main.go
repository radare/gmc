// Copyright 2013 The gmc Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

// A Command is an implementation of a gmc command
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) int

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'gmc help' output.
	Short string

	// Long is the long message shown in the 'gmc help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Usage prints the command usage message and exits with code 2
func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", c.Long)
	os.Exit(1)
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'gmc help'.
var commands = []*Command{
	cmdLs,
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			exitStatus := cmd.Run(cmd, args)
			os.Exit(exitStatus)
		}
	}

	fmt.Fprintf(os.Stderr, "gmc: unknown subcommand %q\nRun 'gmc help' for usage.\n", args[0])
	os.Exit(1)
}

var usageTemplate = `gmc is a full featured commandline Unix like mail client

Usage:
    gmc command [arguments]

The commands are:
{{range .}}    {{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "gmc help [command]" for more information about a command.

`

var helpTemplate = `usage: gmc {{.UsageLine}}
{{.Long}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func usage() {
	tmpl(os.Stderr, usageTemplate, commands)
	os.Exit(1)
}

func help(args []string) {
	if len(args) == 0 {
		usage()
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: gmc help command\n\nToo many arguments given.\n")
		os.Exit(1)
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			os.Exit(1)
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q. Run 'gmc help'.\n", arg)
	os.Exit(1)
}
