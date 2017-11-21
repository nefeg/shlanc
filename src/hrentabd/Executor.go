package hrentabd


type Executor interface {

	Exec(silent bool, jobs ...Job) (outs [][]byte, errs []error)

	OnStart(f ...func(jobs []Job, err error))
	OnComplete(f ...func(jobs []Job, err error))

	OnItemStart(f ...func(job Job, err error, out []byte))
	OnItemComplete(f ...func(job Job, err error, out []byte))
}