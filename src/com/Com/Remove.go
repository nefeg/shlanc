package Com

import (
	"hrentabd"
	"fmt"
	"flag"
	"errors"
	"log"
	"time"
)

type Remove struct{
	Com
}

var MsgUsageComRemove   = "usage: \n\trm -index <index> \n\trm -ts <timestamp> \n\trm --all \n\trm --help\n"

func (c *Remove)Exec(Tab hrentabd.Tab, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = MsgUsageComAdd
		}

	}(&response, &err)

	var INDEX string
	var TS int64
	var ALL, HELP bool

	Args := flag.NewFlagSet("com_remove", flag.PanicOnError)
	Args.StringVar(&INDEX, "index", "", "remove record by index")
	Args.Int64Var(&TS, "ts", 0, "remove records by timestamp")
	Args.BoolVar(&ALL, "all", false, "remove all records")
	Args.BoolVar(&HELP, "help", false, "show this help")
	Args.Parse(args)

	response = MsgUsageComRemove

	if HELP || Args.NFlag() <1 {
		response = MsgUsageComRemove

	}else if ALL{
		Tab.Flush()
		response = "OK"

	}else if INDEX != ""{

		if !Tab.HasJob(INDEX){
			log.Panicln("index not found")
		}

		Tab.RmByIndex(INDEX)
		response = "OK"

	}else if TS != 0 {

		t := time.Unix(TS,0)
		if !Tab.HasJobs(t, true){
			log.Panicf("no jobs found for '%s' \n", t.String())
		}

		response = "OK"
	}

	return response, err
}