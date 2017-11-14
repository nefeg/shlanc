package Com

import (
	"hrentabd"
	"fmt"
	"flag"
	"hrentabd/Job"
	"errors"
	"log"
)

type Add struct{
	Com
}

var MsgUsageComAdd  = "usage: \n\tadd -index <index> -cmd <command to execute> -ttl <ttl> [--force] \n\tadd --help\n"

func (c *Add)Exec(Tab hrentabd.Tab, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = MsgUsageComAdd

			log.Println("[ComAdd]Exec: ", r)
		}

	}(&response, &err)

	var INDEX, CMD string
	var TTL int64
	var OVERRIDE, HELP, REPEAT bool

	Args := flag.NewFlagSet("com_add", flag.PanicOnError)
	Args.StringVar(&INDEX, "index", "", "record index(name/id)? unique string")
	Args.StringVar(&CMD, "cmd", "", "command line for execute")
	Args.Int64Var(&TTL, "ttl", 0, "time to execute (and removing record if -repeat=false)")
	Args.BoolVar(&REPEAT, "repeat", false, "repeat job")
	Args.BoolVar(&OVERRIDE, "force", false, "allow to override existed records")
	Args.BoolVar(&HELP, "help", false, "show this help")
	Args.Parse(args)


	if HELP || Args.NFlag() <1 || (INDEX=="" || CMD=="" || TTL==0){
		response = MsgUsageComAdd

	}else{

		hr := Job.New(INDEX)
		hr.SetCommand(CMD)
		hr.SetTtl(TTL)
		hr.SetRepeatable(REPEAT)

		Tab.PushJobs(OVERRIDE, hr)

		response = "OK"
	}

	return response, err
}