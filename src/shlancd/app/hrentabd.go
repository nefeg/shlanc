package app

import (
	"os"
	"time"
	"fmt"
	. "shlancd/app/api"
	"github.com/umbrella-evgeny-nefedkin/slog"
)

var logPrefix = "[shlancd]"

type hrentabd struct{

	table       Table
	executor    Executor
	options     AppOptions
}


func New(T Table, E Executor, options AppOptions) *hrentabd {

	application := &hrentabd{}
	application.table       = T
	application.executor    = E
	application.options     = options

	return application
}


func (app *hrentabd) Run(){

	defer func(){
		code    := 0
		message := "no message"
		if r:= recover(); r!=nil{
			slog.CritF("%s %v\n", logPrefix, r)
			code = 1
			message = fmt.Sprint(r)
		}

		app.Stop(code, message)
	}()


	app.runHrend(app.options.RunMissed) // todo remove old jobs
}

func (app *hrentabd) Stop(code int, message interface{}){

	app.table.Close()

	slog.InfoF("*** Application terminated with message: %s\n\n", message)

	os.Exit(code)
}

func (app *hrentabd) runHrend(strict bool){

	for{
		currentTime := time.Now()
		if found := app.table.FindByTime(currentTime, strict); len(found)>0{

			go func(list JobListIndex){

				for _, job := range list{

					slog.InfoF("%s Pulling job: %s\n", logPrefix, job.Index())
					if !app.table.PullJob(job){
						slog.InfoF("%s Pulling job: skip (Can't pull job) %s", logPrefix, job.Index())

					}else{

						slog.InfoF("%s Job started: %s", logPrefix, job.Index())
						app.executor.Exec(job)


						if job.IsPeriodic(){
							job.NextPeriod()

							app.table.PushJobs(false, job)
						}
					}
				}

			}(found)
		}

		time.Sleep(1 * time.Second)
	}
}


