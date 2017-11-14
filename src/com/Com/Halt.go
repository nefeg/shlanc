package Com

import (
"hrentabd"
"errors"
	"os"
)

type Halt struct{
	Com
}

var ErrAppHalt = errors.New("** command <HALT> received\n")

func (c *Halt)Exec(Tab hrentabd.Tab, args []string)  (response string, err error){

	print(ErrAppHalt.Error())

	os.Exit(0)

	return "OK", ErrAppHalt
}