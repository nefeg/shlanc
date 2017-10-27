package controls

import (
"hrentabd"
"errors"
	"os"
)

type ComHalt struct{
	Com
}

var ErrAppHalt = errors.New("** command <HALT> received\n")

func (c *ComHalt)Exec(Tab *hrentabd.HrenTab, args []string)  (response string, err error){

	print(ErrAppHalt.Error())

	os.Exit(0)

	return "OK", ErrAppHalt
}