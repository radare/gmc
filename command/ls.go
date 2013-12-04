// Copyright 2013 The gmc Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package command

import (
	"fmt"
)

var cmdLs = &Command{
	UsageLine: "ls [server]",
	Short:     "Short description of the ls command",
	Long: `
Long description of the ls command
`,
	Run: runLs,
}

func runLs(cmd *Command, args []string) error {
	fmt.Println("ls command")
	return nil
}
