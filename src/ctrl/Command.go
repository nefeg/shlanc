package ctrl

import "hrentabd"

type Command interface {

	Resolve(cmdName string) (cmd string, args []string, err error)
	Exec(Tab *hrentabd.HrenTab, args []string)  (string, error)
}

