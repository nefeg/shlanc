package client

import (
	"shlancd/client/telnet"
	"shlancd/client/native"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"shared/config/addr"
)

var logPrefix = "[client.config]"

type Config struct {

	Type    string `json:"type"`

	Options struct {
		Network string `json:"network"`
		Address string `json:"address"`
		Path    string `json:"path"`
	} `json:"options"`
}

func Resolve(conf Config) (client Handler){

	switch conf.Type {
	case "socket":
		client = telnet.New( addr.New(conf.Options.Network, conf.Options.Address) )

	case "native":
		client = native.New()

	default:
		slog.Fatalf("%s Resolve(panic): Unknown client type: %s", logPrefix, conf.Type)
	}

	return client
}