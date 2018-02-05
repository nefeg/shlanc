package client

import (
	"shlancd/cli"
)

type Handler interface {

	Handle(context cli.Context)
}

