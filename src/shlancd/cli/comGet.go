package cli

import (
	"github.com/urfave/cli"
	"errors"
	"fmt"
)

func NewComGet(context *Context) cli.Command {

	return cli.Command{
		Name:    "get",
		Aliases: []string{"g"},
		Usage:   "Get job by id",
		UsageText: "" +
			"\tshlanc get <index>",

		Action:  func(c *cli.Context) (err error) {

			defer func(err *error){
				if r := recover(); r != nil{
					*err = errors.New(fmt.Sprintf("%s", r))
				}
			}(&err)

			jobIndex := c.Args().Get(0)
			if jobIndex == "" {
				panic(ErrCmdArgs)
			}

			c.App.Writer.Write( []byte( view( (*context).Get(jobIndex) ) ) )

			return err
		},
	}
}
