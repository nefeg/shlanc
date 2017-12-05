package hrentabd

import "time"

type Tab interface {
	FindByIndex(index string) Job
	FindByTime(t time.Time, strict bool) IList

	HasJobs(t time.Time, strict bool) bool
	HasJob(index string) bool

	PushJobs(override bool, l ...Job) (pushed int)
	PullJob(job Job) bool

	List() TList
	Flush()

	Close()
}




