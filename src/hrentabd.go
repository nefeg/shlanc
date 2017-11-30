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


	go app.runHrend(false)

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


				app.Exe.OnItemComplete(func(job hrentabd.Job, err error, out []byte){

					fmt.Println("Complete: ", job.Index())

					fmt.Println("==================================================================")
					fmt.Print(string(out))
					fmt.Println("==================================================================")

					// remove executed job if no --repeat flag
					if job.IsRepeatable(){
						job.SetTimeStart(time.Unix(job.TimeStart().Unix() + job.Ttl(),0))

					}else{
						app.Tab.RmByIndex(job.Index())
					}
				})

				app.Exe.Exec(false, arr...)

			}(found)
		}

		time.Sleep(1 * time.Second)
		//print(".")
	}
}


