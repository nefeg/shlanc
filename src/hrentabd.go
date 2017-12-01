package main

import (
	"log"
	"os"
	"time"
	"fmt"
	"hrentabd/Tab"
	"hrentabd"
	"executor"
	"storage"
	"client"
)


type app struct{

	//Api API
	Conf Config

	Tab hrentabd.Tab
	Exe hrentabd.Executor

	Client client.Handler
}


func CreateApp(AppConfig Config) *app {

	application := &app{}
	application.Tab     = hrentabd.Tab( Tab.New( storage.Resolve(AppConfig.Storage) ))
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


	go app.runHrend(false) // todo remove old jobs

	app.Client.Handle(app.Tab)
}

func (app *app) Stop(code int, message interface{}){

	app.Tab.Close()

	log.Printf("*** Application terminated with message: %s\n\n", message)

	os.Exit(code)
}

func (app *app) runHrend(strict bool){

	for{
		currentTime := time.Now()
		if found := app.Tab.FindByTime(currentTime, strict); len(found)>0{

			go func(list hrentabd.IList){

				arr := []hrentabd.Job{}
				for _,v := range list{
					arr = append(arr, v)
				}

				app.Exe.OnItemStart(func(job hrentabd.Job, err error, out []byte){
					log.Println("Job started:", job.Index())
				})

				app.Exe.OnItemComplete(func(job hrentabd.Job, err error, out []byte){
					log.Println("Job complete:", job.Index())

					// remove executed job if no --repeat flag
					if job.IsPeriodic(){
						job.NextPeriod()

						app.Tab.PushJobs(true, job)

					}else{
						app.Tab.RmByIndex(job.Index())
					}
				})

				app.Exe.Exec(arr...)

			}(found)
		}

		time.Sleep(1 * time.Second)
		//print(".")
	}
}


