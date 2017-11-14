package com

import (
	"log"
)

type Config interface{
	Resolve(commandLine string) (cmd Cmd, args []string, err error)
	Add(command Cmd) Config
}

type config struct{
	commands []Cmd
}



func NewConfig(Commands []Cmd) Config{
	return Config( &config{commands:Commands} )
}

func (c *config) Resolve(commandLine string) (cmd Cmd, args []string, err error){

	for _,Com := range c.commands{

		if cmdName, args, err := Com.Resolve(commandLine); err == nil{

			log.Println("Resolved: ", cmdName)

			return Com, args, err
		}
	}

	return cmd, args, err
}


func (c *config) Add(command Cmd) Config{
	c.commands = append(c.commands, command)
	return c
}