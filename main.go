package main

import (
	"fmt"
	"os"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
)

var g_cmd *commander.Commander

func main() {

	g_cmd = &commander.Commander{
		Name: os.Args[0],
		Commands: []*commander.Command{
			git_make_cmd_create(),
			git_make_cmd_login(),
			git_make_cmd_dl_create(),
			git_make_cmd_dl_ls(),
			git_make_cmd_dl_rm(),
		},
		Flag: flag.NewFlagSet("goctogit", flag.ExitOnError),
	}

	err := g_cmd.Flag.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("**error** %v\n", err)
		os.Exit(1)
	}

	args := g_cmd.Flag.Args()
	err = g_cmd.Run(args)
	if err != nil {
		fmt.Printf("**error** %v\n", err)
		os.Exit(1)
	}
}

// EOF
