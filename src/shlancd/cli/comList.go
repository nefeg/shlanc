package cli

import (
	"github.com/urfave/cli"
	"time"
	sapi "shlancd/app/api"
	"errors"
	"fmt"
)

func NewComList(context *Context) cli.Command {

	return cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "Show list of jobs",
		UsageText: "" +
			"\tshlanc list\n" +
			"\tshlanc list -t <timestamp>",

		Flags: 	[]cli.Flag{
			cli.Int64Flag{
				Name:  "timestamp, t",
				Usage: "get jobs at given timestamp",
			},
		},

		Action:  func(c *cli.Context) (err error) {

			defer func(err *error){
				if r := recover(); r != nil{
					*err = errors.New(fmt.Sprintf("%s", r))
				}
			}(&err)


			var jobsList []sapi.Job

			if ts := c.Int64("timestamp"); ts>0{
				jobsList = (*context).ListTime(time.Unix(ts,0))

			}else{
				jobsList = (*context).List()
			}


			response := ""
			for _,job := range jobsList{

				response += view(job)
			}

			c.App.Writer.Write([]byte(response))

			return err
		},
	}
}
