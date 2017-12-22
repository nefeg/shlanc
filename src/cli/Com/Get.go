package Com

import (
	"hrontabd"
	"fmt"
	"flag"
	"errors"
)

type Get struct{
	Com
}

const usage_GET = "usage: \n" +
	"\t  get (\\g) <id>\n" +
	"\t  get (\\g) --help\n"

func (c *Get) Exec(Tab hrontabd.TimeTable, args []string)  (response string, err error){

	defer func(response *string, err *error){
		if r := recover(); r!=nil{
			*err        = errors.New(fmt.Sprint(r))
			*response   = c.Usage()
		}

	}(&response, &err)

	var defaultResponse = "null"

	var HELP, HLP bool

	Args := flag.NewFlagSet("com_get", flag.PanicOnError)
	Args.BoolVar(&HELP, "help", false, "show this help")
	Args.BoolVar(&HLP, "h", false, "show this help")
	Args.Parse(args)

	// show help
	if HELP || HLP || Args.NArg() <1{
		response = c.Usage()

		//show item by index
	}else {

		if found := Tab.FindJob(Args.Arg(0)); found != nil{
			response += c.view(found)
		}
	}

	if response == "" {
		response = defaultResponse
	}

	return response, nil
}

func (c *Get) Usage() string{

	return c.Desc() + "\n\t" + usage_GET
}

func (c *Get) view(job hrontabd.Job) string{

	return fmt.Sprintln(
		job.Id(),"\t",
		job.CronLine(),"\t",
		"\""+job.Command()+"\"", "\t",
		"\""+job.Comment()+"\"", "\t",
	)
}

