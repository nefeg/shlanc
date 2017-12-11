package Tab

import (
	"time"
	. "hrentabd"
	j "hrentabd/Job"
	"log"
)

type tab struct {

	db      Storage
	hTList  TList    // time list
	version string
}


func New (s Storage) *tab{

	t := &tab{ db:s, hTList:TList{} }

	t.load()

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

		log.Printf("[Tab]PushJobs: %d jobs pushed", *pushed)

		if r := recover(); r!=nil{
			log.Printf("[Tab]PushJobs (panic): %v", r)
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

func (h *tab) PullJob(job Job) bool{

	if jobString := h.db.Pull(job.Index()); jobString != ""{
		job.UnSerialize(jobString)
		h.rmByIndex(job.Index())
		return true
	}

	return false
}




func (h *tab) List() TList {

	h.sync()

	return h.hTList
}

func (h *tab) Flush() {

	h.hTList = TList{}
	h.db.Flush()
}

func (h *tab) Close(){
	h.db.Disconnect()
}



func (h *tab) sync(){

	if h.version != h.db.Version(){
		log.Printf("sync: %v --> %v\n", h.version, h.db.Version())
		h.load()
	}
}

func (h *tab) load() (loaded int){

	h.version   = h.db.Version()
	h.hTList    = TList{}
	for index,data := range h.db.List(){
		h.loadJob(j.New(string(index)).UnSerialize(string(data)))
		loaded++
	}

	return loaded
}

func (h *tab) loadJob(job Job){

	currentTs := job.TimeStart().Unix()
	if _,isset := h.hTList[currentTs]; !isset{
		h.hTList[currentTs] = IList{}
	}

	h.hTList[currentTs][job.Index()] = job
}


// remove jobs from local index (not from storage)
func (h *tab) rmByIndex(index string) bool{

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

func (h *tab) rmByTime(t time.Time, strict bool) bool{

	if list := h.FindByTime(t, strict); len(list)>0{

		for i := range list{
			h.rmByIndex(i)
		}

		return true
	}

	return false
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
