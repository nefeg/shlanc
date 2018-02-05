package cli

import (
	sapi "shlancd/app/api"
	"time"
)

type Context interface{

	List() []sapi.Job
	ListTime(tm time.Time) []sapi.Job

	Get(index string) sapi.Job

	Add(job sapi.Job, force bool) bool

	Remove(id string)
	RemoveTime(tm time.Time)

	Purge()
	Term()
}