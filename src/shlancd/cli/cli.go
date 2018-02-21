package cli

import (
	"time"
	"github.com/urfave/cli"
)

func New() *cli.App{

	Cli := cli.NewApp()
	Cli.Version             = "0.4"
	Cli.Name                = "ShLANC-cli"
	Cli.Usage               = "[SH]lanc [L]ike [A]s [N]ot [C]ron"
	Cli.Author              = "Evgeny Nefedkin"
	Cli.Compiled            = time.Now()
	Cli.Email               = "evgeny.nefedkin@umbrella-web.com"
	Cli.Description         = "Distributed and concurrency manager of deffer jobs"

	return Cli
}
