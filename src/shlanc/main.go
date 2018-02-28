package main

import (
	"github.com/urfave/cli"
	"fmt"
	"errors"
	"shared/sig"
	. "shared/config/app"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"os"
	"shlancd/storage"

	sapi "shlancd/app/api"
	capi "shlancd/cli"
	capiCtx "shlancd/cli/Context"
)

var ConfigPaths         = []string{
	"config.json",
	"/etc/shlanc/config.json",
	"/etc/shlacd/config.json",
}


var ErrCmdArgs          = errors.New("ERR: expected argument")

const FL_DEBUG      = "debug"

var sigIntHandler = func(){}

func init()  {

	slog.SetLevel(slog.LvlNone)
	slog.SetFormat(slog.FormatTimed)

	sig.SIG_INT(&sigIntHandler)

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "client version=%s {%s}\n", c.App.Version, c.App.Compiled.Format("2006/01/02 15:04:05"))
	}
}


func main(){

	var CliContext capi.Context

	sigIntHandler = func(){
		CliContext.Term()
	}

	defer func(a *capi.Context){
		if r := recover(); r != nil{

			slog.Debugln("[main ->defer] ", r)

			if r == ErrCmdArgs || r == ErrNoConfFile{
				fmt.Println(r)
				fmt.Println("See: `shlanc --help` or `shlanc <command> --help`")

			}else if r == ErrConfCorrupted{
				fmt.Println(r)

			}else{
				fmt.Println(r)
			}
		}
		if *a != nil{
			(*a).Term()
		}

	}(&CliContext)


	Cli := capi.New()

	Cli.Before = func(c *cli.Context) error {

		// debug flag
		if c.GlobalBool(FL_DEBUG){
			slog.SetLevel(slog.LvlDebug)
			slog.Debugln("[main] os.Args: ", os.Args)
		}

		// Override config
		if confFile := c.GlobalString("config"); confFile != ""{
			ConfigPaths = []string{confFile}
		}
		slog.Debugln("[main ->Cli.Before] Config paths: ", ConfigPaths)


		mainConfig := LoadConfig(ConfigPaths)

		JTable := sapi.NewTable( storage.Resolve(mainConfig.Storage) )

		slog.Debugln("[main ->Cli.Before] JTable: ", JTable)

		CliContext = capiCtx.New(JTable)

		slog.Debugln("[main ->Cli.Before] CliContext: ", CliContext)

		return nil
	}


	// CONFIG
	Cli.Flags =  []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "path to daemon config-file",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "show debug log",
		},
	}

	// COMMANDS
	Cli.Commands = []cli.Command{
		capi.NewComAdd(&CliContext,),
		capi.NewComList(&CliContext),
		capi.NewComRemove(&CliContext),
		capi.NewComPurge(&CliContext),
		capi.NewComGet(&CliContext),
	}

	Cli.Run(os.Args)
}

