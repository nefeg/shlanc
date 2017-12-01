package hrentabd

import "time"

type Job interface {
	Ttl()               int64
	Index()             string
	Command()           string
	TimeStart()         time.Time

	SetTtl(ttl int64)
	SetCommand(command string)
	SetTimeStart(t time.Time)

	// repeatable
	IsPeriodic() bool
	SetPeriod(period int64)
	GetPeriod()(period int64)
	NextPeriod()


	Serialize() string
	UnSerialize(data string) Job
}
