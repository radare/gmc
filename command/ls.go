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
