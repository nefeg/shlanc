package main

import (
	"runtime"
	"sig"
	"fmt"
	"log"
	"executor"
	. "com"
	. "com/Com"
	"client-api/telnet"
	"client-api"
	"hrentabd"
	"storage"
)


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

	App := CreateApp(
		hrentabd.Executor( executor.NewExecutorLocal() ),
		//hrentabd.Storage( storage.NewStorageFile("/tmp/hren.db") ),
		hrentabd.Storage( storage.NewStorageRedis("tcp", "127.0.0.1:6379") ),
	)

	go App.Run()

	telnetHandler := telnet.NewHandler( telnet.NewConnectionConf("tcp", "127.0.0.1:6607") )
	client_api.Handler(telnetHandler).Handle(App.Tab, CommandConfig)
}