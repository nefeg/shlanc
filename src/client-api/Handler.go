package client_api

import (
	"hrentabd"
	"com"
)

type Handler interface {

	Handle(Tab hrentabd.Tab, Conf com.Config)
}

