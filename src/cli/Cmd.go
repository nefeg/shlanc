package cli

import "hrentabd"

type Cmd interface {

	Resolve(cmdName string) (cmd string, args []string, err error)
	Exec(Tab hrentabd.Tab, args []string)  (string, error)
	Usage() string
}