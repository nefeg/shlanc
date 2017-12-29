package client

import (
	"shlacd/hrontabd"
)

type Handler interface {

	Handle(Tab hrontabd.TimeTable)
}

