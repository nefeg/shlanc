package executor

import (
	"hrentabd"
	"os/exec"
	"sync"
)

type localExecutor struct {

	silent  bool
	async   bool

	onStart         []func(jobs []hrentabd.Job, err error)
	onComplete      []func(jobs []hrentabd.Job, err error)

	onItemStart     []func(job hrentabd.Job, err error, out []byte)
	onItemComplete  []func(job hrentabd.Job, err error, out []byte)
}


func NewExecutorLocal(silent, async bool) *localExecutor{

	return &localExecutor{ silent:silent, async:async }
}


// execute list of jobs
func (a *localExecutor)Exec(jobs ...hrentabd.Job) (outs [][]byte, errs []error){

	var err error
	var out []byte

	// on start
	for _,f := range a.onStart {
		f(jobs, err)
	}


	wg := &sync.WaitGroup{}
	run := func(job hrentabd.Job, wg *sync.WaitGroup){
		wg.Add(1)

		out, err = a.ExecItem(job, wg)
		errs    = append(errs, err)
		outs    = append(outs, out)

		wg.Done()
	}

	for _,job := range jobs{

		if a.async{
			go run(job,wg)
		}else{
			run(job,wg)
		}
	}


	wg.Wait()

	// on complete
	for _,f := range a.onComplete {
		f(jobs, err)
	}

	return outs, errs
}


func (a *localExecutor)ExecItem(job hrentabd.Job, wg *sync.WaitGroup) (out []byte, err error){

	// on item start
	for _,f := range a.onItemStart {
		f(job, err, out)
	}

	// RUN COMMAND
	if cmd := exec.Command("sh",  "-c", job.Command()); !a.silent {
		out,err = cmd.Output()
	}

	// on item complete
	for _,f := range a.onItemComplete{
		f(job, err, out)
	}

	return out, err
}


func (a *localExecutor)OnStart(f ...func(jobs []hrentabd.Job, err error)){ a.onStart = f }
func (a *localExecutor)OnComplete(f ...func(jobs []hrentabd.Job, err error)){ a.onComplete = f }

func (a *localExecutor)OnItemStart(f ...func(job hrentabd.Job, err error, out []byte)){ a.onItemStart = f }
func (a *localExecutor)OnItemComplete(f ...func(job hrentabd.Job, err error, out []byte)){ a.onItemComplete = f }