package main

import (
	"runtime"
	"shared/sig"
	. "shared/config/app"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"github.com/urfave/cli"
	"time"
	"os"

	"shlancd/app/api"
	"shlancd/storage"
	"shlancd/app"
	"shlancd/executor"
	"shlancd/client"
	"shlancd/cli/Context"
	"fmt"
)

var sigIntHandler   = func(){}
var logPrefix       = "[main]"
var ConfigPaths     = []string{
	"config.json",
	"/etc/shlanc/config.json",
	"/etc/shlancd/config.json",
}

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())

	sig.SIG_INT(&sigIntHandler)

	slog.SetLevel(slog.LvlInfo)
	slog.SetFormat(slog.FormatTimed)
}


func main(){

	const FL_DEBUG  = "debug"
	const FL_CONFIG = "config"

	var application Application
	var table       api.Table
	var config      *Config

	Cli := cli.NewApp()
	Cli.Version             = "0.23"
	Cli.Name                = "ShLANC-server"
	Cli.Usage               = "[SH]lac [L]ike [A]s [N]ot [C]ron"
	Cli.Author              = "Evgeny Nefedkin"
	Cli.Compiled            = time.Now()
	Cli.Email               = "evgeny.nefedkin@umbrella-web.com"
	Cli.EnableBashCompletion= true
	Cli.Description         = "Distributed and concurrency manager of deffer jobs"

	Cli.Before = func(c *cli.Context) error {

		if c.GlobalBool(FL_DEBUG){
			slog.SetLevel(slog.LvlDebug)
			slog.DebugF("%s Starting...\n", logPrefix)
			slog.DebugLn(logPrefix, " Args:", os.Args)
		}

		// Override config
		if confFile := c.GlobalString(FL_CONFIG); confFile != ""{
			ConfigPaths = []string{confFile}
		}
		slog.DebugLn(logPrefix, " Config paths:", ConfigPaths)

		config = LoadConfig(ConfigPaths)
		//
		appOptions := api.AppOptions{RunMissed: config.RunMissed}

		table = api.NewTable( storage.Resolve(config.Storage) )
		//
		application = app.New(
			table,
			executor.Resolve(config.Executor),
			appOptions,
		)
		//
		sigIntHandler = func(){
			application.Stop(1, sig.ErrSigINT)
		}

		return nil
	}

	 //CONFIG
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

	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s {%s}\n", c.App.Version, c.App.Compiled.Format("2006/01/02 15:04:05"))
	}

	Cli.Action = func(c *cli.Context) {

		go application.Run()

		client.Resolve(config.Client).Handle( Context.New(table) )
	}

	Cli.Run(os.Args)
}