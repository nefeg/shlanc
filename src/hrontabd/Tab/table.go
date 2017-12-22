package Tab

import (
	. "hrontabd"
	"log"
	J "hrontabd/Job"
)

type table struct {

	db          Storage
	jobs        []Job
	version     string
}


// constructor
func New (s Storage) *table{

	t := &table{ db:s }
	// t.load()

	return t
}



func (t *table) FindJob(jobId string) (job Job){

	// t.sync()
	for _, j := range t.jobs{
		log.Println(j.Id(),jobId)
		if j.Id() == jobId { return j }
	}

	return nil
}

func (t *table) RmJob(jobId string) bool{

	for i := range t.jobs{
		if t.jobs[i].Id() == jobId{
			t.jobs = append(t.jobs[:i], t.jobs[i+1:]...)
			return true
		}
	}

	//r := t.db.Rm(jobId)
	//t.sync()

	return false
}

func (t *table) AddJob(job Job, force bool){

	defer func(job Job){

		if r := recover(); r!=nil{
			log.Printf("[Tab]Add (panic): %v", r)
			panic(r)
		}

		log.Printf("[Tab]Add: job#%s added", job.Id())

	}(job)


	if t.FindJob(job.Id()) != nil && !force{
		panic("job already exist")
	}


	t.jobs = append(t.jobs, job)
}

func (t *table) PullJob(jobId string) (job Job){

	log.Println("[hrentab.table] PullJob: Trying to lock job...")
	//if t.db.Lock(jobId){
	//
	//	log.Printf("[hrentab.table] PullJob: Job #%s locked\n", jobId)
	//	if jobData := t.db.Get(jobId); jobData != ""{
	//		job = J.New().UnSerialize(string(jobData))
	//	}
	//}else{
	//	log.Printf("[hrentab.table] PullJob: Locking for job#%s fail\n", jobId)
	//}

	for _,j := range t.jobs{
		if j.Id() == jobId{
			job = j
		}
	}

	return job
}

func (t *table) PushJob(job Job)  {

	t.db.UnLock(job.Id())
}

func (t *table) ListJobs() []Job{

	t.sync()

	return t.jobs
}

func (t *table) Flush() {

	t.db.Flush()
	t.sync()
}

func (t *table) Close(){
	t.db.Disconnect()
}



func (t *table) sync(){

	if !t.isSynced(){
		log.Printf("sync: %v --> %v\n", t.version, t.db.Version())
		t.load()
	}
}

func (t *table) isSynced() bool{
	return t.version == t.db.Version()
}

func (t *table) load(){

	for _, jobData := range t.db.List() {
		t.jobs = append(t.jobs, J.New().UnSerialize(string(jobData)) )
	}
}