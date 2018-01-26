package client

import (
	"hrentabd/app/api"
)

type Handler interface {

	Handle(Tab api.Table)
}

