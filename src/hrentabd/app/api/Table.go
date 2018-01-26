package api

import (
	"time"
	"log"
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

	t := &table{ db:s, hTList:JobListTime{} }

	t.load()

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

		log.Printf("[JobTable]PushJobs: %d jobs pushed", *pushed)

		if r := recover(); r!=nil{
			log.Printf("[JobTable]PushJobs (panic): %v", r)
			panic(r)
		}

	}(&pushed)

	for _,job := range l{

		if h.HasJob(job.Index()){
			if override{
				h.PullJob(job)
			}else{
				panic("Job index already exist")
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
		log.Printf("sync: %v --> %v\n", h.version, h.db.Version())
		h.load()
	}
}

func (h *table) load() (loaded int){

	h.version   = h.db.Version()
	h.hTList    = JobListTime{}
	for index,data := range h.db.List(){
		h.loadJob(NewJob(string(index)).UnSerialize(string(data)))
		loaded++
	}

	return loaded
}

func (h *table) loadJob(job Job){

	currentTs := job.TimeStart().Unix()
	if _,isset := h.hTList[currentTs]; !isset{
		h.hTList[currentTs] = JobListIndex{}
	}

	h.hTList[currentTs][job.Index()] = job
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


