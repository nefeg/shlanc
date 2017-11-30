package executor

import (
	"log"
	"hrentabd"
)

// executor config
type Config struct {
	Type    string `json:"type"`

	Options struct {

	} `json:"options"`
}

func Resolve(conf Config) (exe hrentabd.Executor){

	switch conf.Type {
	case "local":
		exe = NewExecutorLocal()

	default:
		log.Panicln("Unknown client type")
	}

	return exe
}