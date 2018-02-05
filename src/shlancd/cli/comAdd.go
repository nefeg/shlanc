package cli

import (
	"github.com/urfave/cli"
	"github.com/satori/go.uuid"

	sapi "shlancd/app/api"
	"errors"
	"time"
	"fmt"
)

func NewComAdd(context *Context) cli.Command {

	var DateFormat      = "2006-01-02 15:04:05 -07"

	return cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add job",
		UsageText: "" +
			"\tshlanc add [-i <index>] [-r <seconds>] [--force] -e <command to execute> -l <ttl> \n" +
			"\tshlanc add [-i <index>] [-r <seconds>] [--force] -e <command to execute> -s <timestamp> \n" +
			"\tshlanc add [-i <index>] [-r <seconds>] [--force] -e <command to execute> -t <" + DateFormat + "> \n",

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "cmd,e",
				Usage: "Command to execute (required)",
			},
			cli.StringFlag{
				Name:  "index,i",
				Usage: "Set job index manually",
			},
			cli.Int64Flag{
				Name:  "ttl, l",
				Usage: "Set job ttl (pause before execute). Required one of: ttl, ts or tm",
			},
			cli.Int64Flag{
				Name:  "ts, s",
				Usage: "Set time at job will start (timestamp format). Required one of: ttl, ts or tm",
			},
			cli.StringFlag{
				Name:  "tm, t",
				Usage: `Set time at job to execute (format: "` + DateFormat + `"). Required one of: ttl, ts or tm`,
			},
			cli.Int64Flag{
				Name: "repeat, r",
				Usage: "Repeat job every N sec (N>0, the time is after the first run, " +
					"e.g. -ttl=30 -repeat=10 - run job after 30sec(ttl) and then every 10sec)" +
					"If flag is NOT set, job will be deleted automatically after it executed.",
			},
			cli.BoolFlag{
				Name:  "force",
				Usage: "Override duplicates",
			},
		},


		Action: func(c *cli.Context) (err error) {

			defer func(err *error){
				if r := recover(); r != nil{
					*err = errors.New(fmt.Sprintf("%s", r))
				}
			}(&err)

			CMD := c.String("cmd")
			if CMD == "" {
				panic(errors.New(`ERR: "-cmd" required` + "\nsee `add --help`"))
			}

			TTL, TS, TM := c.Int64("ttl"), c.Int64("ts"), c.String("tm")
			if TTL == 0 && TS == 0 && TM == "" {
				panic(errors.New("ERR: Expected option, one of -ttl, -tm or -ts \nsee `add --help`"))
			}

			INDEX := c.String("index")
			if INDEX == "" {
				INDEX = uuid.NewV4().String()
			}

			REPEAT := c.Int64("repeat")
			OVERRIDE := c.Bool("force")

			newJob := sapi.NewJob(INDEX)
			newJob.CommandX(CMD)
			newJob.PeriodX(REPEAT)

			if TTL > 0 {
				newJob.TtlX(TTL)

			} else if TS > 0 {
				newJob.AtX(time.Unix(TS, 0))

			} else if TM != "" {
				if t, err := time.Parse(DateFormat, TM); err == nil {
					newJob.AtX(t)

				} else {
					panic(err)
				}
			}

			if (*context).Add(newJob, OVERRIDE) {
				c.App.Writer.Write( []byte(newJob.Index()) )
			}

			return err
		},
	}
}
