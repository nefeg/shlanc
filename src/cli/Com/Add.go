package Com

import (
	"hrentabd"
	"fmt"
	"flag"
	"hrentabd/Job"
	"errors"
	"log"
	"time"
	"github.com/satori/go.uuid"
)

type Add struct{
	Com
}

const usage_ADD  = "usage: \n" +
	"\t  add (\\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -ttl <ttl> \n" +
	"\t  add (\\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -ts <timestamp> \n" +
	"\t  add (\\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -tm <2006-01-02T15:04:05Z07:00> \n" +
	"\t  add (\\a) --help\n"

const dateFormat = "2006-01-02 15:04:05 -0700 MST"

func (c *Add)Exec(Tab hrentabd.Tab, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()

			log.Println("[ComAdd]Exec: ", r)
		}

	}(&response, &err)

	var INDEX, CMD, TM string
	var TTL, TS, REPEAT int64
	var OVERRIDE, HELP, HLP bool

	Args := flag.NewFlagSet("com_add", flag.PanicOnError)
	Args.StringVar(&INDEX,      "index",    "",     "record index(name/id)? unique string")
	Args.StringVar(&CMD,        "cmd",      "",     "command line for execute")
	Args.StringVar(&TM,         "tm",       "",     "time (format: "+dateFormat+") to start and removing")
	Args.Int64Var(&TTL,         "ttl",      0,      "time (seconds) to start (and removing record if -repeat=false)")
	Args.Int64Var(&TS,          "ts",       0,      "timestamp (seconds) at start and removing record")
	Args.Int64Var(&REPEAT,      "repeat",   0,      "repeat period")
	Args.BoolVar(&OVERRIDE,     "force",    false,  "allow to override existed records")
	Args.BoolVar(&HELP,         "help",     false,  "show this help")
	Args.BoolVar(&HLP,          "h",        false,  "show this help")
	Args.Parse(args)


	if HELP || HLP || Args.NFlag() <1 || (CMD=="" || (TTL==0 && TS==0 && TM=="")){
		response = c.Usage()

	}else{

		if INDEX==""{
			INDEX = uuid.NewV4().String()
		}

		hr := Job.New(INDEX)
		hr.SetCommand(CMD)
		hr.SetPeriod(REPEAT)

		if TTL>0 {
			hr.SetTtl(TTL)

		}else if TS>0 {
			hr.SetTimeStart(time.Unix(TS,0))

		}else if TM!=""{
			if t, err := time.Parse(dateFormat, TM); err ==nil{
				hr.SetTimeStart(t)

			}else {
				response = err.Error()
			}
		}


		Tab.PushJobs(OVERRIDE, hr)

		response = "OK"
	}

	return response, err
}

func (c *Add) Usage() string{
	return c.Desc() + "\n\t" + usage_ADD
}