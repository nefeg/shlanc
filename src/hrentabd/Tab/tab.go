package Tab

import (
	"storage"
	"time"
	. "hrentabd"
	j "hrentabd/Job"
	"log"
)

type tab struct {

	db      storage.Storage
	hTList  TList    // time list
}


func New (s storage.Storage) *tab{

	t := &tab{ db:s, hTList:TList{} }

	for index,data := range s.List(){
		t.loadJob(j.New(string(index)).UnSerialize(string(data)))
	}

	return t
}

// interface Tab
func (h *tab) FindByIndex(index string) Job{

	for _,jobs := range h.List(){
		if job, isset := jobs[index]; isset{
			return job
		}
	}

	return nil
}

func (h *tab) FindByTime(t time.Time, strict bool) IList {

	var list IList

	ts := t.Unix()
	for k, hmap := range h.List() {
		if (strict && k == ts) || (!strict && k <= ts) {
			list = mergeIndexList(list, hmap)
		}
	}

	return list
}


func (h *tab) RmByIndex(index string) bool{

	for ts,jobs := range h.List(){
		if job, isset := jobs[index]; isset{

			delete(h.hTList[ts], index)

			h.db.Rm(storage.Index(job.Index()))

			// remove empty ts-key
			if tm := time.Unix(ts,0); !h.HasJobs(tm, true){
				delete(h.hTList, ts)
			}

			return true
		}
	}

	return false
}

func (h *tab) RmByTime(t time.Time, strict bool) bool{

	if list := h.FindByTime(t, strict); len(list)>0{

		for i := range list{
			h.RmByIndex(i)
		}

		return true
	}

	return false
}


// check jobs by time start
func (h *tab) HasJobs(t time.Time, strict bool) bool {

	ts := t.Unix()
	for k, hList := range h.List() {
		if (!strict && k <= ts || strict && k == ts) && len(hList) > 0 {
			return true
		}
	}
	return false
}

// check job by index
func (h *tab) HasJob(index string) bool{

	return h.FindByIndex(index) != nil
}



func (h *tab) PushJobs(override bool, l ...Job) (pushed int){

	defer func(pushed *int){

		log.Printf("[tab]PushJobs: %d jobs pushed", *pushed)

		if r := recover(); r!=nil{
			log.Panicln("[tab]PushJobs: ", r)
		}

	}(&pushed)

	for _,job := range l{

		if h.HasJob(job.Index()){
			if override{
				h.RmByIndex(job.Index())
			}else{
				log.Panicln("Job index already exist")
			}
		}

		h.loadJob(job)

		h.db.Add(
			storage.Index(job.Index()),
			storage.Record(job.Serialize()),
			override)

		pushed++
	}

	return pushed
}

func (h *tab) List() TList {
	return h.hTList
}

func (h *tab) Flush() {

	h.hTList = TList{}
	h.db.Flush()
}


func (h *tab) loadJob(job Job){

	currentTs := job.TimeStart().Unix()
	if _,isset := h.hTList[currentTs]; !isset{
		h.hTList[currentTs] = IList{}
	}

	h.hTList[currentTs][job.Index()] = job
}

func mergeIndexList(m ...IList) IList{

	var r   = IList{}

	for _, a := range m{
		for k, v := range a {
			r[k] = Job(v)
		}
	}

	return r
}
