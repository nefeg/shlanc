package main

import (
	"time"
	"hrentabd"
	"runtime"
	"fmt"
	"sig"
	"ctrl"
	"os"
	"executer"
)

const VERSION       = "0.1"
const CTRL_PROTO    = "tcp"
//const CTRL_PROTO    = "unix"
const CTRL_ADDR     = "127.0.0.1:6607"
//const CTRL_ADDR     = "/tmp/hren.sock"

var IPCChan     = make(chan string)
var logChan     = make(chan string)
var runStrict   = false
var Tab         = &hrentabd.HrenTab{}


func init()  {

	runtime.GOMAXPROCS(runtime.NumCPU())

	sig.SIG_INT(func(){
		fmt.Println("Callback")
	})
}


func main(){
	println("Starting...")

	defer func(){
		if r:= recover(); r!=nil{
			os.Exit(1)
		}
	}()

	runLogger()

	runHrend(Tab, false)

	ctrl.Run( ctrl.NewCommandHandler(Tab), ctrl.NewConnectionConf(CTRL_PROTO, CTRL_ADDR))
}


func runHrend(Tab *hrentabd.HrenTab, strict bool){

	go func(){

		e := executer.NewAsyncExecutor()

		for{
			currentTime := time.Now()
			if found := Tab.FindByTime(currentTime, strict); len(found)>0{

				go func(list hrentabd.IndexList){

					arr := []hrentabd.Hren{}
					for _,v := range list{
						arr = append(arr, v)
					}


					e.OnItemComplete(func(job hrentabd.Hren, err error, out []byte){

						fmt.Println("Complete: ", job.Index())

						fmt.Println("==================================================================")
						fmt.Print(string(out))
						fmt.Println("==================================================================")

						// remove executed job if no --repeat flag
						if job.IsRepeatable(){
							job.SetTimeStart(time.Unix(job.TimeStart().Unix() + job.Ttl(),0))

						}else{
							Tab.RmByIndex(job.Index())
						}
					})

					e.Exec(false, arr...)

				}(found)
			}

			time.Sleep(1 * time.Second)
			//print(".")
		}
	}()
}

func runLogger(){

	go func(ch chan string){

		for {println(<-ch)}

	}(logChan)
}