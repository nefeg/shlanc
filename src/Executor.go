package main

import "hrentabd"

type Executor interface {

	Exec(silent bool, jobs ...hrentabd.Hren) (outs [][]byte, errs []error)

	OnStart(f ...func(jobs []hrentabd.Hren, err error))
	OnComplete(f ...func(jobs []hrentabd.Hren, err error))

	OnItemStart(f ...func(job hrentabd.Hren, err error, out []byte))
	OnItemComplete(f ...func(job hrentabd.Hren, err error, out []byte))
}