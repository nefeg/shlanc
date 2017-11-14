package executor

import "hrentabd"

type Executor interface {

	Exec(silent bool, jobs ...hrentabd.Job) (outs [][]byte, errs []error)

	OnStart(f ...func(jobs []hrentabd.Job, err error))
	OnComplete(f ...func(jobs []hrentabd.Job, err error))

	OnItemStart(f ...func(job hrentabd.Job, err error, out []byte))
	OnItemComplete(f ...func(job hrentabd.Job, err error, out []byte))
}