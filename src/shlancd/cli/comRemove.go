package cli

import (
	"github.com/urfave/cli"
	"time"
	"errors"
	"fmt"
)

func NewComRemove(context *Context) cli.Command {

	return cli.Command{
		Name:    "remove",
		Aliases: []string{"rm", "r"},
		Usage:   "Remove jobs by ID or time of start ",
		UsageText: "" +
			"\tshlanc remove <job id>\n" +
			"\tshlanc remove --all\n" +
			"\tshlanc remove -t <timestamp>",

		Flags: 	[]cli.Flag{
			cli.BoolFlag{
				Name:  "all,purge",
				Usage: "remove all jobs",
			},
			cli.Int64Flag{
				Name:  "timestamp,t",
				Usage: "remove jobs with given time",
			},
		},

		Action:  func(c *cli.Context) (err error) {

			defer func(err *error){
				if r := recover(); r != nil{
					*err = errors.New(fmt.Sprintf("%s", r))
				}
			}(&err)

			if ts := c.Int64("timestamp"); ts > 0 {
				(*context).RemoveTime(time.Unix(ts,0))

			}else if c.Bool("all"){
				(*context).Purge()

			}else if jobId := c.Args().Get(0); jobId != ""{
				(*context).Remove(jobId)

			}else{
				panic(ErrCmdArgs)
			}

			return err
		},
	}
}