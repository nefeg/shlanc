package main

import (
	"log"
	"os"
	"time"
	"fmt"
	"hrentabd/Tab"
	"hrentabd"
)

type app struct{

	//Api Api

	Tab hrentabd.Tab
	Exe hrentabd.Executor
}


func CreateApp(ex hrentabd.Executor, db hrentabd.Storage) *app{

	application := &app{}
	application.Tab = hrentabd.Tab( Tab.New( db ))
	application.Exe = ex

	return application
}


func (app *app) Run(){

	defer func(){
		code := 0
		if r:= recover(); r!=nil{
			log.Println(r)
			code = 1
		}

		app.Exit(code)
	}()


	app.runHrend(false)
}

func (app *app) Exit(code int){

	app.Tab.Close()

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


