package executer

import (
	"hrentabd"
	"os/exec"
)

type asyncExecutor struct {
	onStart         []func(jobs []hrentabd.Hren, err error)
	onComplete      []func(jobs []hrentabd.Hren, err error)

	onItemStart     []func(job hrentabd.Hren, err error, out []byte)
	onItemComplete  []func(job hrentabd.Hren, err error, out []byte)
}


func NewAsyncExecutor() *asyncExecutor{

	return &asyncExecutor{}
}


// execute list of jobs
func (a *asyncExecutor)Exec(silent bool, jobs ...hrentabd.Hren) (outs [][]byte, errs []error){

	var err error
	var out []byte

	// on start
	for _,f := range a.onStart {
		f(jobs, err)
	}

	for _,job := range jobs{

		out, err = a.ExecItem(silent, job)

		errs    = append(errs, err)
		outs    = append(outs, out)
	}


	// on complete
	for _,f := range a.onComplete {
		f(jobs, err)
	}

	return outs, errs
}


func (a *asyncExecutor)ExecItem(silent bool, job hrentabd.Hren) (out []byte, err error){

	// on item start
	for _,f := range a.onItemStart {
		f(job, err, out)
	}

	// RUN COMMAND
	if cmd := exec.Command("sh",  "-c", job.Command()); !silent {
		out,err = cmd.Output()
	}

	// on item complete
	for _,f := range a.onItemComplete{
		f(job, err, out)
	}

	return out, err
}


func (a *asyncExecutor)OnStart(f ...func(jobs []hrentabd.Hren, err error)){ a.onStart = f }
func (a *asyncExecutor)OnComplete(f ...func(jobs []hrentabd.Hren, err error)){ a.onComplete = f }

func (a *asyncExecutor)OnItemStart(f ...func(job hrentabd.Hren, err error, out []byte)){ a.onItemStart = f }
func (a *asyncExecutor)OnItemComplete(f ...func(job hrentabd.Hren, err error, out []byte)){ a.onItemComplete = f }