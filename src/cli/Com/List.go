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
	"\t  list (\\l) -i <index>\n" +
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
			response = fmt.Sprintln("==> ", found.TimeStart().String(), "(", found.TimeStart().Unix() ,")")
			response += fmt.Sprintln(INDEX, ":", found.Command(), "(", found.Ttl() ,")")
		}


	// show items by timestamp
	}else if TS != 0{

		if found := Tab.FindByTime(time.Unix(TS,0), false); found != nil{

			for ind, job := range found{

				response += fmt.Sprintln(
					job.TimeStart().String(), "(", job.TimeStart().Unix() ,")","\t",
					ind,  "\t",
					job.Command(), "\t",
					job.Ttl(),
				)
			}
		}

	// show all jobs
	}else{

		for ts, ah := range Tab.List(){

			for index, h := range ah{
				response += fmt.Sprintln(
					time.Unix(ts,0).String(), "(", ts ,")","\t",
					index,  "\t",
					h.Command(), "\t",
					h.Ttl(),
				)
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

