package app

import (
	"shlancd/storage"
	"shlancd/client"
	"shlancd/executor"
)

type Config struct {

	// run missed jobs on start
	RunMissed   bool

	// storage config
	Storage  storage.Config `json:"storage"`

	// client config
	Client client.Config `json:"client"`

	// executor config
	Executor executor.Config `json:"executor"`
}
