package api

import (
	"time"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

//###############################################
//                  INTERFACE
//###############################################
type Table interface {
	FindByIndex(index string) Job
	FindByTime(t time.Time, strict bool) JobListIndex

	HasJobs(t time.Time, strict bool) bool
	HasJob(index string) bool

	PushJobs(override bool, l ...Job) (pushed int)
	PullJob(job Job) bool

	List() JobListTime
	Flush()

	Close()
}


//###############################################
//                  CONSTRUCTOR
//###############################################
func NewTable (s Storage) *table{

	slog.Debugln("[api.Table] NewTable: creating...")

	t := &table{ db:s, hTList:JobListTime{} }

	t.load()

	slog.Debugln("[api.Table] NewTable: ", t)

	return t
}



//###############################################
//                IMPLEMENTATION
//###############################################
type table struct {

	db      Storage
	hTList  JobListTime    // time list
	version string
}


// interface JobTable
func (h *table) FindByIndex(index string) Job{

	for _,jobs := range h.List(){
		if job, isset := jobs[index]; isset{
			return job
		}
	}

	return nil
}

func (h *table) FindByTime(t time.Time, strict bool) JobListIndex {

	var list = JobListIndex{}

	ts := t.Unix()
	for k, hmap := range h.List() {
		if (strict && k == ts) || (!strict && k <= ts) {
			list.Merge(hmap)
		}
	}

	return list
}

// check jobs by time start
func (h *table) HasJobs(t time.Time, strict bool) bool {

	ts := t.Unix()
	for k, hList := range h.List() {
		if (!strict && k <= ts || strict && k == ts) && len(hList) > 0 {
			return true
		}
	}
	return false
}


// check job by index
func (h *table) HasJob(index string) bool{

	return h.FindByIndex(index) != nil
}



func (h *table) PushJobs(override bool, l ...Job) (pushed int){

	defer func(pushed *int){

		slog.Infof("[api.Table] PushJobs: %d jobs pushed\n", *pushed)

		if r := recover(); r!=nil{
			slog.Panicf("[api.Table] PushJobs (panic): %v\n", r)
			panic(r)
		}

	}(&pushed)

	for _,job := range l{

		if h.HasJob(job.Index()){
			if override{
				h.PullJob(job)
			}else{
				panic("job index already exist")
			}
		}

		h.loadJob(job)

		h.db.Push(
			job.Index(),
			job.Serialize(),
			override)

		pushed++
	}

	return pushed
}



func (h *table) PullJob(job Job) bool{

	slog.Debugln("[api.Table] PullJob: ", job)

	if jobString := h.db.Pull(job.Index()); jobString != ""{
		job.UnSerialize(jobString)
		h.rmByIndex(job.Index())
		return true
	}

	return false
}


func (h *table) List() JobListTime {

	h.sync()

	return h.hTList
}



func (h *table) Flush() {

	h.hTList = JobListTime{}
	h.db.Flush()
}



func (h *table) Close(){
	h.db.Disconnect()
}





func (h *table) sync(){

	if h.version != h.db.Version(){
		slog.Debugf("[api.Table] sync: %v --> %v\n", h.version, h.db.Version())
		h.load()
	}
}

func (h *table) load() (loaded int){

	slog.Debugln("[api.Table] _load: loading database")

	h.version   = h.db.Version()
	h.hTList    = JobListTime{}
	for index,data := range h.db.List(){
		h.loadJob(NewJob(string(index)).UnSerialize(string(data)))
		loaded++
	}

	return loaded
}

func (h *table) loadJob(job Job){

	slog.Debugln("[api.Table] _loadJob: ", job)

	currentTs := job.At().Unix()
	if _,isset := h.hTList[currentTs]; !isset{
		h.hTList[currentTs] = JobListIndex{}
	}

	h.hTList[currentTs][job.Index()] = job

	slog.Debugln("[api.Table] _loadJob: hTList --> ", job)
}


// remove jobs from local index (not from storage)
func (h *table) rmByIndex(index string) bool{

	for ts,jobs := range h.List(){
		if _, isset := jobs[index]; isset{

			delete(h.hTList[ts], index)

			// remove empty ts-key
			if tm := time.Unix(ts,0); !h.HasJobs(tm, true){
				delete(h.hTList, ts)
			}

			return true
		}
	}

	return false
}

func (h *table) rmByTime(t time.Time, strict bool) bool{

	if list := h.FindByTime(t, strict); len(list)>0{

		for i := range list{
			h.rmByIndex(i)
		}

		return true
	}

	return false
}


