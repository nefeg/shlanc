package storage

import (
	"shlancd/app/api"
	"shlancd/storage/redis"
	"shlancd/storage/file"
	"github.com/umbrella-evgeny-nefedkin/slog"
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

func Resolve(conf Config) (storage api.Storage){

	switch conf.Type {
	case "redis":
		slog.Debugln("[storage.config] Resolve: redis")
		storage = redis.New(conf.Options.Network, conf.Options.Address, conf.Options.Key)

	case "file":
		slog.Debugln("[storage.config] Resolve: file")
		storage = file.New(conf.Options.Path)

	case "script":
		slog.Debugln("[storage.config] Resolve: script")
		slog.Panicln("[storage.config] Resolve: script-type is not implemented yet")
		// todo implement this

	default:
		slog.Panicln("[storage.config] Resolve: unknown storage type")
	}

	return storage
}
