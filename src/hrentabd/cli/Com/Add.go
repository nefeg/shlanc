package Com

import (
	"fmt"
	"flag"
	"hrentabd/app/api"
	"errors"
	"time"
	"github.com/satori/go.uuid"
	"log"
)

type Add struct{
	Com
}

const usage_ADD  = "usage: \n" +
	"\t  add (\\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -ttl <ttl> \n" +
	"\t  add (\\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -ts <timestamp> \n" +
	"\t  add (\\a) [-index <index>] [--force] [--repeat <seconds>] -cmd <command to execute> -tm <2006-01-02T15:04:05Z07:00> \n" +
	"\t  add (\\a) --help\n"

const dateFormat = "2006-01-02 15:04:05 -07"

func (c *Add)Exec(Tab api.Table, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New("ERR: " + fmt.Sprint(r))
			*response   = c.Usage()

			log.Println(r)
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

		newJob := api.NewJob(INDEX)
		newJob.SetCommand(CMD)
		newJob.SetPeriod(REPEAT)

		if TTL>0 {
			newJob.SetTtl(TTL)

		}else if TS>0 {
			newJob.SetTimeStart(time.Unix(TS,0))

		}else if TM!=""{
			if t, err := time.Parse(dateFormat, TM); err ==nil{
				newJob.SetTimeStart(t)

			}else {
				panic(err)
			}
		}


		Tab.PushJobs(OVERRIDE, newJob)

		response = "OK"
	}

	return response, err
}

func (c *Add) Usage() string{
	return c.Desc() + "\n\t" + usage_ADD
}