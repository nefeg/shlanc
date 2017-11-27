package Com

import (
	"hrentabd"
	"errors"
)

type Quit struct{
	Com
}

const usage_QUIT = "usage: \n\t  quit (\\q) \n"

var ErrConnectionClosed = errors.New("** command <QUIT> received")

func (c *Quit)Exec(Tab hrentabd.Tab, args []string)  (string, error){

	return "OK", ErrConnectionClosed
}

func (c *Quit) Usage() string{
	return c.Desc() + "\n\t" + usage_QUIT
}