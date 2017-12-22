package cli

import "hrontabd"

type Cmd interface {

	Resolve(cmdName string) (cmd string, args []string, err error)
	Exec(Tab hrontabd.TimeTable, args []string)  (string, error)
	Usage() string
}