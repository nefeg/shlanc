package hrentabd

import "time"

type Hren interface {
	Ttl()               int64
	Index()             string
	Command()           string
	TimeStart()         time.Time
	IsRepeatable()      bool

	SetTtl(ttl int64)
	SetTimeStart(t time.Time)
	SetRepeatable(repeat bool)
}
