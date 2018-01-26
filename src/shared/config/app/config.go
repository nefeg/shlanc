package app

import (
	"hrentabd/storage"
	"hrentabd/client"
	"hrentabd/executor"
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
