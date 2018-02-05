package cli

import (
	"github.com/urfave/cli"
)

func NewComPurge(context *Context) cli.Command {

	return cli.Command{
		Name:    "purge",
		Usage:   "Remove all jobs ",
		UsageText: "" +
			"shlanc purge",

		Action:  func(c *cli.Context) error {

			(*context).Purge()

			return nil
		},
	}
}
