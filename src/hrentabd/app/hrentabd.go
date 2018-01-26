package app

import (
	"log"
	"os"
	"time"
	"fmt"
	"hrentabd/executor"
	"hrentabd/storage"
	"hrentabd/client"

	. "shared/config/app"
	. "hrentabd/app/api"
)


type hrentabd struct{

	//Api API
	Conf Config

	Table Table
	Exe   Executor

	Client client.Handler
}


func New(AppConfig Config) *hrentabd {

	application := &hrentabd{}
	application.Table   = NewTable( storage.Resolve(AppConfig.Storage) )
	application.Exe     = executor.Resolve(AppConfig.Executor)
	application.Client  = client.Resolve(AppConfig.Client)

	return application
}


func (app *hrentabd) Run(){

	defer func(){
		code    := 0
		message := "no message"
		if r:= recover(); r!=nil{
			log.Println(r)
			code = 1
			message = fmt.Sprint(r)
		}

		app.Stop(code, message)
	}()


	go app.runHrend(app.Conf.RunMissed) // todo remove old jobs

	app.Client.Handle(app.Table)
}

func (app *hrentabd) Stop(code int, message interface{}){

	app.Table.Close()

	log.Printf("*** Application terminated with message: %s\n\n", message)

	os.Exit(code)
}

func (app *hrentabd) runHrend(strict bool){

	for{
		currentTime := time.Now()
		if found := app.Table.FindByTime(currentTime, strict); len(found)>0{

			go func(list JobListIndex){

				for _, job := range list{

					log.Println("[SYSTEM] Pulling job:", job.Index())
					if !app.Table.PullJob(job){
						log.Println("[SYSTEM] Pulling job: skip (Can't pull job)", job.Index())

					}else{

						log.Println("[SYSTEM] Job started:", job.Index())
						app.Exe.Exec(job)


						if job.IsPeriodic(){
							job.NextPeriod()

							app.Table.PushJobs(false, job)
						}
					}
				}

			}(found)
		}

		time.Sleep(1 * time.Second)
	}
}


