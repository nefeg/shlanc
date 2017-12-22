package hrontabd

import "time"

type Job interface {
	Id()                string
	Command()           string
	CronLine()          string
	Comment()           string
	TimeStart()         time.Time

	SetCommand(command string)
	SetCronLine(timeLine string)
	SetComment(comment string)

	Serialize() string
	UnSerialize(data string) Job
}
