package client

import (
	"hrentabd"
)

type Handler interface {

	Handle(Tab hrentabd.Tab)
}

