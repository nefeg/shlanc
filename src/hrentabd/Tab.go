package hrentabd

import "time"

type Tab interface {
	FindByIndex(index string) Job
	FindByTime(t time.Time, strict bool) IList

	RmByIndex(index string) bool
	RmByTime(t time.Time, strict bool) bool

	HasJobs(t time.Time, strict bool) bool
	HasJob(index string) bool

	PushJobs(override bool, l ...Job) (pushed int)

	List() TList
	Flush()
}




