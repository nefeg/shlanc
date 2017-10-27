package hrentabd

import (
	"time"
)

type IndexList map[string]Hren
type TimeList map[int64]IndexList


type HrenTab struct {
	// todo remove hren
	hTList  TimeList    // time list
}

func mergeIndexList(m ...IndexList) IndexList{

	var r IndexList = IndexList{}

	for _, a := range m{
		for k, v := range a {
			r[k] = Hren(v)
		}
	}

	return r
}

func (h *HrenTab)this() *HrenTab{
	if h.hTList == nil {
		h.hTList    = TimeList{}
	}
	return h
}

func (h *HrenTab)FindByIndex(index string) Hren{

	for _,jobs := range h.List(){
		if job, isset := jobs[index]; isset{
			return job
		}
	}

	return nil
}

func (h *HrenTab) FindByTime(t time.Time, strict bool) IndexList {

	var list IndexList

	ts := t.Unix()
	for k, hmap := range h.List() {
		if (strict && k == ts) || (!strict && k <= ts) {
			list = mergeIndexList(list, hmap)
		}
	}

	return list
}


func (h *HrenTab)RmByIndex(index string) bool{

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

func (h *HrenTab)RmByTime(t time.Time, strict bool) bool{

	if list := h.FindByTime(t, strict); len(list)>0{

		for i := range list{
			h.RmByIndex(i)
		}

		return true
	}

	return false
}


// check jobs by time start
func (h *HrenTab)HasJobs(t time.Time, strict bool) bool {

	ts := t.Unix()
	for k, hList := range h.List() {
		if (!strict && k <= ts || strict && k == ts) && len(hList) > 0 {
			return true
		}
	}
	return false
}

// check job by index
func (h *HrenTab) HasJob(index string) bool{

	return h.FindByIndex(index) != nil
}

func (h *HrenTab) PushJobs(override bool, l ...Hren) *HrenTab {

	if len(l) == 0 {
		panic("No jobs for push")
	}

	for _,job := range l{

		if h.HasJob(job.Index()){
			if override{
				h.RmByIndex(job.Index())
			}else{
				panic("Job index already exist")
			}
		}


		currentTs := job.TimeStart().Unix()
		if _,isset := h.this().hTList[currentTs]; !isset{
			h.this().hTList[currentTs] = IndexList{}
		}

		h.this().hTList[currentTs][job.Index()] = job

		//// add job to time-group
		//if h.HasJobs(job.TimeStart(), true) {
		//
		//	currentTs := job.TimeStart().Unix()
		//	h.this().hTList[currentTs][job.Index()] = job // currentTimeIList is type of IndexList
		//
		//// add job to NEW time-group
		//}else{
		//	h.this().hTList[job.TimeStart().Unix()] = IndexList{job.Index():job}
		//}
	}

	return h
}

func (h *HrenTab) List() TimeList {
	return h.this().hTList
}

func (h *HrenTab) Flush() {

	h.hTList = TimeList{}
}
