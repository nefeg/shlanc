package controls

import (
	"hrentabd"
	"errors"
)

type ComQuit struct{
	Com
}

var ErrConnectionClosed = errors.New("** command <QUIT> received")

func (c *ComQuit)Exec(Tab *hrentabd.HrenTab, args []string)  (string, error){

	return "OK", ErrConnectionClosed
}