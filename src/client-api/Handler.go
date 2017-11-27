package client_api

import (
	"hrentabd"
	"cli"
)

type Handler interface {

	Handle(Tab hrentabd.Tab, Conf cli.CLI)
}

