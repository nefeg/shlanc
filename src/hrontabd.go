package main

import (
	"log"
	"os"
	"time"
	"fmt"
	"hrontabd/Tab"
	"hrontabd"
	"executor"
	"storage"
	"client"
)


type app struct{

	//Api API
	Conf Config

	Tab hrontabd.TimeTable
	Exe hrontabd.Executor

	Client client.Handler
}


func CreateApp(AppConfig Config) *app {

	application := &app{}
	application.Tab     = hrontabd.TimeTable( Tab.New( storage.Resolve(AppConfig.Storage) ))
	application.Exe     = executor.Resolve(AppConfig.Executor)
	application.Client  = client.Resolve(AppConfig.Client)

	return application
}


func (app *app) Run(){

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

	app.Client.Handle(app.Tab)
}

func (app *app) Stop(code int, message interface{}){

	app.Tab.Close()

	log.Printf("*** Application terminated with message: %s\n\n", message)

	os.Exit(code)
}

func (app *app) runHrend(strict bool){

	for{
		if found := app.Tab.ListJobs(); len(found)>0{

			go func(jobs []hrontabd.Job){

				for _, job := range jobs{

					log.Println("Pulling job:", job.Id())
					if j := app.Tab.PullJob(job.Id()); j != nil{

						log.Println("Job started:", j.Id())
						//app.Exe.Exec(job)
						//app.Tab.PushJob(job)

					}else{
						log.Println("Pulling job: skip (Can't pull job)", job.Id())

					}
				}

			}(found)
		}

		time.Sleep(1 * time.Minute)
	}
}


