package Com

import (
	"errors"
	"strings"
	"github.com/mattn/go-shellwords"
)

type Com struct {
	name            string
	alias           string
	desc            string
	resolvedName    string

}

var ErrComUnresolved    = errors.New("")

func New(name, alias, desc string) Com{

	return Com{name:name, alias:alias, desc:desc}
}

func (c *Com)Resolve(cmdLine string) (cmd string, args []string, err error){

	parts := strings.Fields(cmdLine)


	if len(parts) > 0 && ( c.resolveAs(parts[0], c.Name()) || c.resolveAs(parts[0], c.Alias())) {
		cmd         = c.Name()
		args,err    = shellwords.Parse( strings.Join(parts[1:]," ") )

	}else{
		err = ErrComUnresolved
	}

	return cmd, args, err
}

func (c *Com)resolveAs(cmdLine, resolveName string) (resolved bool){

	if cmdLine == resolveName {
		c.resolvedName  = resolveName
		resolved        = true
	}

	return resolved
}

func (c *Com)resolvedAs() string{
	return c.resolvedName
}


func (c *Com)Name() string{
	return c.name
}

func (c *Com)Alias() string{
	return c.alias
}

func (c *Com)Desc() string{
	return c.desc
}