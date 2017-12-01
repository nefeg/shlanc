package Com

import (
	"hrentabd"
	"fmt"
	"time"
	"flag"
	"errors"
	"log"
)

type List struct{
	Com
}

const usage_LIST = "usage: \n" +
	"\t  list (\\l) -index <index>\n" +
	"\t  list (\\l) -ts <timestamp> \n" +
	"\t  list (\\l) --help\n"

func (c *List) Exec(Tab hrentabd.Tab, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()
		}

	}(&response, &err)

	var defaultResponse = "(empty)"

	var INDEX string
	var TS int64
	var HELP, HLP bool

	Args := flag.NewFlagSet("com_list", flag.PanicOnError)
	Args.StringVar(&INDEX, "index", "", "search by index")
	Args.Int64Var(&TS, "ts", 0, "search by timestamp")
	Args.BoolVar(&HELP, "help", false, "show this help")
	Args.BoolVar(&HLP, "h", false, "show this help")
	Args.Parse(args)

	// Args.PrintDefaults()


	if INDEX != "" && TS != 0{
		log.Panicln("[ComList] Exec: ", ErrComTooMuchArgs)
	}

	// show help
	if HELP || HLP{
		response = c.Usage()

	//show item by index
	}else if INDEX != ""{

		if found := Tab.FindByIndex(INDEX); found != nil{
			response = c.view(found)
		}


	// show items by timestamp
	}else if TS != 0{

		if found := Tab.FindByTime(time.Unix(TS,0), false); found != nil{
			for _, job := range found{
				response += c.view(job)
			}
		}

	// show all jobs
	}else{

		for _, jobsGroup := range Tab.List(){
			for _, job := range jobsGroup{
				response += c.view(job)
			}
		}
	}


	if response == "" {
		response = defaultResponse
	}

	return response, nil
}

func (c *List) Usage() string{

	return c.Desc() + "\n\t" + usage_LIST
}

func (c *List) view(job hrentabd.Job) string{

	var period string
	if job.IsPeriodic(){
		period = fmt.Sprint(job.GetPeriod())
	}else{
		period = "null"
	}

	return fmt.Sprintln(
		job.TimeStart().String(),"\t",
		job.Index(),"\t",
		period,"\t",
		"\""+job.Command()+"\"", "\t",
		job.Ttl(),"\t",
	)
}

