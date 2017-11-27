package telnet

import (
	"hrentabd"
	"net"
	"io"
	"fmt"
	"log"
	"os"
	"cli"
	"cli/Com"
)

type handler struct {

	addr    net.Addr
}

const WlcMessage = "HrenTab terminal connected OK"


func NewHandler(listen net.Addr) *handler{

	return &handler{addr:listen}
}


func (h *handler) Handle(Tab hrentabd.Tab, Cli cli.CLI){

	IPC, err := net.Listen(h.addr.Network(), h.addr.String())
	if err != nil {
		log.Panicf("%s: %s", "ERROR", err.Error())
	}
	defer func(){
		IPC.Close()
		if UAddr, err := net.ResolveUnixAddr(h.addr.Network(), h.addr.String()); err == nil{
			os.Remove(UAddr.String())
		}
	}()

	for{
		if Connection, err := IPC.Accept(); err == nil {

			h.handleConnection(Connection, Tab, Cli)
			Connection.Close()

		}else{
			log.Println(err.Error())
			continue
		}
	}
}

func (h *handler)handleConnection(Connection net.Conn, Tab hrentabd.Tab, Cli cli.CLI){

	var response string

	defer func(response *string){

		if r := recover(); r != nil{

			if r == io.EOF {
				*response = "client socket closed."
				writeData(Connection, "\n"+(*response)+"\n")
				log.Println("\nSession closed by cause: " + (*response))

			}else{

				writeData(Connection, "\n" + fmt.Sprint(r) + "\n")
				panic(r)
			}
		}else{
			writeData(Connection, "\n" + (*response) + "\n")
			log.Println("\nSession closed by cause: " + (*response))
		}
	}(&response)


	//writeData(Connection, WlcMessage + "\n")
	writeData(Connection, WlcMessage + "\n>>")
	for{

		if rcv, err := readData(Connection); rcv != ""{

			if err != nil {
				log.Println(err.Error())
				response = err.Error()

			}else{

				response = "Unknown command"
				if Command, args, err := Cli.Resolve(rcv); Command != nil{

					response, err = Command.Exec(Tab, args)

					if err != nil{
						response = err.Error()

						// ComQuit
						if err == Com.ErrConnectionClosed{
							return
						}
					}
				}else{
					response += "\n" + Cli.Help()
				}
			}

			writeData(Connection,                                                                                                                                                                                                                                                                                                                                                                                                                                                                       response+ "\n>>")
			//writeData(Connection, "==> " + response + "\n>>")
		}
	}
}



