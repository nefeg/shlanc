package Com

import (
	"hrontabd"
	"fmt"
	"flag"
	"hrontabd/Job"
	"errors"
	"log"
	"github.com/satori/go.uuid"
)

type Add struct{
	Com
}

const usage_ADD  =
	"\\a -cron <cron-line> -cmd <command> [-id <id>] [-comment <comment>] [--force]\n" +
	"\t\\a --help\n"


func (c *Add)Exec(Tab hrontabd.TimeTable, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()

			log.Println("[ComAdd]Exec: ", r)
		}

	}(&response, &err)

	var INDEX, CMD, COMMENT, CLINE string
	var OVERRIDE, HELP, HLP bool

	Args := flag.NewFlagSet("com_add", flag.PanicOnError)
	Args.StringVar(&INDEX,      "id",       "",     "record index(name/id)? unique string")
	Args.StringVar(&CLINE,      "cron",     "",     "cron-formatted time line")
	Args.StringVar(&CMD,        "cmd",      "",     "command")
	Args.StringVar(&COMMENT,    "comment",  "",     "comment")
	Args.BoolVar(&OVERRIDE,     "force",    false,  "allow to override existed records")
	Args.BoolVar(&HELP,         "help",     false,  "show this help")
	Args.BoolVar(&HLP,          "h",        false,  "show this help")
	Args.Parse(args)


	if HELP || HLP || CMD=="" || CLINE == ""{
		response = c.Usage()

	}else{

		if INDEX==""{
			INDEX = uuid.NewV4().String()
		}

		log.Println(COMMENT)

		job := Job.New()
		job.SetID(INDEX)
		job.SetCronLine(CLINE)
		job.SetCommand(CMD)
		job.SetComment(COMMENT)

		log.Println(job)

		Tab.AddJob(job, OVERRIDE)

		response = "OK"
	}


	return response, err
}

func (c *Add) Usage() string{
	return c.Desc() + "\n\t" + usage_ADD
}