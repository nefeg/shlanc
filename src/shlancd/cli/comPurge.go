package cli

import (
	"github.com/urfave/cli"
)

func NewComPurge(context *Context) cli.Command {

	return cli.Command{
		Name:    "purge",
		Usage:   "remove all jobs ",
		UsageText: "Example: " +
			"hren-cli2 purge",

		Action:  func(c *cli.Context) error {

			(*context).Purge()

			return nil
		},
	}
}
