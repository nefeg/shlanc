package cli

import (
	"github.com/urfave/cli"
	"time"
)

func NewComRemove(context *Context) cli.Command {

	return cli.Command{
		Name:    "remove",
		Aliases: []string{"rm", "r"},
		Usage:   "remove jobs ",
		UsageText: "Example: \n",

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

		Action:  func(c *cli.Context) error {

			if ts := c.Int64("timestamp"); ts > 0 {
				(*context).RemoveTime(time.Unix(ts,0))

			}else if c.Bool("all"){
				(*context).Purge()

			}else if jobId := c.Args().Get(0); jobId != ""{
				(*context).Remove(jobId)

			}else{
				panic(ErrCmdArgs)
			}

			return nil
		},
	}
}