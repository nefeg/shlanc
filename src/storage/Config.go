package storage

import (
	"log"
	"hrentabd"
	"storage/redis"
	"storage/file"
)

// storage config
type Config struct {
	Type    string `json:"type"`

	Options struct {
		Network string `json:"network"`
		Address string `json:"address"`
		Key     string `json:"key"`
		Path    string `json:"path"`
	} `json:"options"`
}

func Resolve(conf Config) (storage hrentabd.Storage){

	switch conf.Type {
	case "redis":
		storage = redis.New(conf.Options.Network, conf.Options.Address, conf.Options.Key)

	case "file":
		storage = file.New(conf.Options.Path)

	case "script":
		// todo implement this

	default:
		log.Panicln("Unknown storage type")
	}

	return storage
}