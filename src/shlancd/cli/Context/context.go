package Context

import "time"
import (
	sapi "shlancd/app/api"
)


func New(T sapi.Table) *context{


	return &context{T}
}

type context struct {

	table   sapi.Table
}

func (c *context) List() []sapi.Job {

	t := c.table.List()
	i := t.ToIndexList()

	return i.ToArray()
}

func (c *context) ListTime(tm time.Time) []sapi.Job{

	li := c.table.FindByTime(tm, true)

	return li.ToArray()

}


func (c *context) Get(index string) sapi.Job{

	return c.table.FindByIndex(index)
}

func (c *context) Add(job sapi.Job, force bool) bool{

	return c.table.PushJobs(force, job) > 0
}



func (c *context) Remove(index string){

	if job := c.table.FindByIndex(index); job != nil{
		c.table.PullJob(job)
	}
}

func (c *context) RemoveTime(tm time.Time){

	if jobs := c.table.FindByTime(tm, true); len(jobs) >0 {
		for _, job := range jobs{
			c.table.PullJob(job)
		}
	}
}

func (c *context) Purge(){
	c.table.Flush()
}

func (c *context) Term(){
	c.table.Close()
}


