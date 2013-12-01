// Copyright 2013 The gmc Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/radare/gmc/command"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	if args[0] == "help" {
		help(args[1:])
		os.Exit(1)
	}

	for _, cmd := range command.Commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() ; os.Exit(1) }
			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()
			err := cmd.Run(cmd, args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			} else {
				return
			}
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
	tmpl(os.Stderr, usageTemplate, command.Commands)
}

func help(args []string) {
	if len(args) == 0 {
		usage()
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: gmc help command\n\nToo many arguments given.\n")
		return
	}

	arg := args[0]

	for _, cmd := range command.Commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q. Run 'gmc help'.\n", arg)
}
