package cli

import "hrentabd/app/api"

type Cmd interface {

	Resolve(cmdName string) (cmd string, args []string, err error)
	Exec(Tab api.Table, args []string)  (string, error)
	Usage() string
}