package Com

import (
	"hrentabd"
	"errors"
)

type Quit struct{
	Com
}

var ErrConnectionClosed = errors.New("** command <QUIT> received")

func (c *Quit)Exec(Tab hrentabd.Tab, args []string)  (string, error){

	return "OK", ErrConnectionClosed
}