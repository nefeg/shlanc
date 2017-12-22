package client

import (
	"hrontabd"
)

type Handler interface {

	Handle(Tab hrontabd.TimeTable)
}

