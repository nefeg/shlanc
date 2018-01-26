package cli

import (
	"log"
)

type CLI interface{
	Resolve(commandLine string) (cmd Cmd, args []string, err error)
	Add(command Cmd) CLI
	List() []Cmd
	Help() (help string)
}




type cli struct{
	commands    []Cmd
}

func New() CLI{
	return CLI( &cli{ commands:commandConfig } )
}

func (c *cli) Resolve(commandLine string) (cmd Cmd, args []string, err error){

	for _,Com := range c.commands{

		if cmdName, args, err := Com.Resolve(commandLine); err == nil{

			log.Println("[CLI] Resolved: ", cmdName)

			return Com, args, err
		}
	}

	return cmd, args, err
}

func (c *cli) Add(command Cmd) CLI{
	c.commands = append(c.commands, command)
	return c
}

func (c *cli) List() []Cmd{
	return c.commands
}

func (c *cli) Help() (help string){
	for _,cmd := range c.commands{
		help += cmd.Usage() +"\n"
	}

	return help
}