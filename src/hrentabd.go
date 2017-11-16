package main

import (
	"time"
	"hrentabd"
	"runtime"
	"fmt"
	"sig"
	"client-api"
	"os"
	LocalExec "executor/Local"
	//"storage/file"
	"storage"
	"hrentabd/Tab"
	"client-api/telnet"
	. "com"
	. "com/Com"
	"log"
	"storage/redis"
)

// const VERSION       = "0.2"
const CTRL_PROTO    = "tcp"
//const CTRL_PROTO    = "unix"
const CTRL_ADDR     = "127.0.0.1:6607"
//const CTRL_ADDR     = "/tmp/hren.sock"

var DB  storage.Storage
var T   hrentabd.Tab

var CommandConfig = NewConfig(
	[]Cmd{

		Cmd( &Halt{New("halt",    `\t`)} ),
		Cmd( &Quit{New("exit",    `\q`)} ),
		Cmd( &List{New("list",    `\l`)} ),
		Cmd( &Add{New("add",    `\a`)} ),
		Cmd( &Remove{New("rm",  `\r`)} ),
	})

func init()  {

	runtime.GOMAXPROCS(runtime.NumCPU())

	sig.SIG_INT(func(){
		fmt.Println("Callback")
	})
}


func main(){
	log.Println("Starting...")

	DB  = storage.Storage( redis.NewRedisStorage("tcp", "127.0.0.1:6379"))
	//DB  = storage.Storage( file.NewFileStorage("/tmp/hren.db") )
	T   = hrentabd.Tab( Tab.New(DB))

	defer func(){
		if r:= recover(); r!=nil{
			os.Exit(1)
		}

		DB.Disconnect()
	}()

	runHrend(T, false)

	telnetHandler := telnet.NewHandler( telnet.NewConnectionConf(CTRL_PROTO, CTRL_ADDR) )
	client_api.Handler(telnetHandler).Handle(T, CommandConfig)
}


func runHrend(Tab hrentabd.Tab, strict bool){

	go func(){

		e := LocalExec.New()

		for{
			currentTime := time.Now()
			if found := Tab.FindByTime(currentTime, strict); len(found)>0{

				go func(list hrentabd.IList){

					arr := []hrentabd.Job{}
					for _,v := range list{
						arr = append(arr, v)
					}


					e.OnItemComplete(func(job hrentabd.Job, err error, out []byte){

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