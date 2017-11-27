package main

import (
	"runtime"
	"log"
	"executor"
	"cli"
	"client-api/telnet"
	"client-api"
	"hrentabd"
	"storage"
)



func init()  {

	runtime.GOMAXPROCS(runtime.NumCPU())

	//sig.SIG_INT(func(){
	//	fmt.Print("Terminateing application...")
	//	panic(sig.ErrSigINT)
	//	fmt.Println("OK")
	//})
}


func main(){
	log.Println("Starting...")

	Application := CreateApp(
		hrentabd.Executor( executor.NewExecutorLocal() ),
		//hrentabd.Storage( storage.NewStorageFile("/tmp/hren.db") ),
		hrentabd.Storage( storage.NewStorageRedis("tcp", "127.0.0.1:6379") ),
	)

	go Application.Run()
	//defer func() {
	//	if r := recover(); r==sig.ErrSigINT{
	//		Application.Exit(0)
	//	}
	//}()

	telnetHandler := telnet.NewHandler( telnet.NewConnectionConf("tcp", "127.0.0.1:6607") )
	client_api.Handler(telnetHandler).Handle(Application.Tab, cli.New())
}