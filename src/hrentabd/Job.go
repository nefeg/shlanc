package hrentabd

import "time"

type Job interface {
	Ttl()               int64
	Index()             string
	Command()           string
	TimeStart()         time.Time
	IsRepeatable()      bool

	SetTtl(ttl int64)
	SetCommand(command string)
	SetTimeStart(t time.Time)
	SetRepeatable(repeat bool)

	Serialize() string
	UnSerialize(data string) Job
}
