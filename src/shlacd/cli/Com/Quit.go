package Com

import (
	"shlacd/hrontabd"
	"errors"
)

type Quit struct{
	Com
}

const usageQuit = "usage: \n\t  quit (\\q) \n"

var ErrConnectionClosed = errors.New("** command <QUIT> received")

func (c *Quit)Exec(Tab hrontabd.TimeTable, args []string)  (string, error){

	return "OK", ErrConnectionClosed
}

func (c *Quit) Usage() string{
	return c.Desc() + "\n\t" + usageQuit
}