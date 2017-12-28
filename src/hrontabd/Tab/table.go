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
	t.load()

	return t
}



func (t *table) FindJob(jobId string) (job Job){

	t.sync()
	for _, j := range t.jobs{
		if j.Id() == jobId { return j }
	}

	return nil
}

func (t *table) RmJob(jobId string) bool{

	r := t.db.Rm(jobId)
	t.sync()

	return r
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

	if !t.db.Lock(job.Id()){
		panic("can't get lock")
	}

	defer t.db.UnLock(job.Id())

	t.db.Add(job.Id(), job.Serialize(), force)

	t.sync()
	//t.jobs = append(t.jobs, job)

}

func (t *table) PullJob(jobId string) (job Job){

	log.Println("[hrentab.table] PullJob: Trying to lock job...")
	if t.db.Lock(jobId){

		log.Printf("[hrentab.table] PullJob: Job #%s locked\n", jobId)
		job = t.FindJob(jobId)
		if job == nil{
			t.db.UnLock(jobId)
			log.Printf("[hrentab.table] PullJob: Job '%s' not found\n", jobId)
		}

	}else{
		log.Printf("[hrentab.table] PullJob: Locking for job#%s fail\n", jobId)
	}

	return job
}

func (t *table) PushJob(job Job)  {

	t.db.UnLock(job.Id())
	log.Printf("[hrentab.table] PushJob: Release lock for job#%s\n", job.Id())

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

	t.jobs    = nil
	t.version = t.db.Version()
	for _, jobData := range t.db.List() {
		t.jobs = append(t.jobs, J.New().UnSerialize(string(jobData)) )
	}
}