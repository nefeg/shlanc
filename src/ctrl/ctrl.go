package ctrl

import (
	"net"
	"log"
	"os"
	"fmt"
)



func Start(handler CommandHandler, listen net.Addr){

	go Run(handler, listen)
}

func Run(handler CommandHandler, listen net.Addr){

	IPC, err := net.Listen(listen.Network(), listen.String())
	if err != nil {
		log.Panicf("%s: %s", "ERROR", err.Error())
	}
	defer func(){
		IPC.Close()
		if UAddr, err := net.ResolveUnixAddr(listen.Network(), listen.String()); err == nil{
			os.Remove(UAddr.String())
		}
	}()

	for{
		if Connection, err := IPC.Accept(); err == nil {

			handler.Handle(Connection)
			Connection.Close()

		}else{
			fmt.Println(err.Error())
			continue
		}
	}
}
