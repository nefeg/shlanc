package native

import (
	"shlancd/cli"
)

type handler struct{}

func New() *handler{

	return &handler{}
}


func (h *handler) Handle(context cli.Context){

	// nothing here to do
}