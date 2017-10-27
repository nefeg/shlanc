package Hren

import (
	"hrentabd"
	"time"
)

type hren struct {
	index       string
	command     string
	startAt     time.Time
	repeat      bool
}

func New(index, command string) hrentabd.Hren{
	return hrentabd.Hren( &hren{index:index, command:command} )
}

func (h *hren)Ttl() int64{
	return h.startAt.Unix() - time.Now().Unix()
}

func (h *hren)SetTtl(ttl int64){
	h.startAt = time.Unix(time.Now().Unix() + int64(ttl), 0)
}


func (h *hren)Index() string{
	return h.index
}

func (h *hren)Command() string{
	return h.command
}

// time start
func (h *hren)TimeStart() time.Time{
	return h.startAt
}

func (h *hren)SetTimeStart(t time.Time){
	h.startAt = t
}

// Repeatable
func (h *hren)IsRepeatable() bool{
	return h.repeat
}

func (h *hren)SetRepeatable(repeat bool){
	h.repeat = repeat
}