package socket

import (
	"hrontabd"
	"net"
	"io"
	"fmt"
	"log"
	"os"
	ComLineIf "cli"
	"cli/Com"
)

type handler struct {

	addr    net.Addr
	cli     ComLineIf.CLI
}

const WlcMessage = "ShLAC terminal connected OK\n" +
	"type \"help\" or \"\\h\" for show available commands"


func NewHandler(listen net.Addr, cli ComLineIf.CLI) *handler{

	return &handler{ addr:listen, cli:cli }
}


func (h *handler) Handle(Tab hrontabd.TimeTable){

	IPC, err := net.Listen(h.addr.Network(), h.addr.String())
	if err != nil {
		log.Panicf("%s: %s", "ERROR", err.Error())
	}
	log.Println("[cient.socket] Listen:", IPC.Addr().String())

	defer func(){
		IPC.Close()
		if UAddr, err := net.ResolveUnixAddr(h.addr.Network(), h.addr.String()); err == nil{
			os.Remove(UAddr.String())
		}
	}()

	for{
		if Connection, err := IPC.Accept(); err == nil {

			go func(){
				h.handleConnection(Connection, Tab)
				Connection.Close()
			}()

		}else{
			log.Println(err.Error())
			continue
		}
	}
}

func (h *handler)handleConnection(Connection net.Conn, Tab hrontabd.TimeTable){

	var response string

	defer func(response *string){

		if r := recover(); r != nil{

			if r == io.EOF {
				*response = "client socket closed."
				writeData(Connection, "\n"+(*response)+"\n")
				log.Println("Session closed by cause: " + (*response))

			}else{

				writeData(Connection, "\n" + fmt.Sprint(r) + "\n")
				panic(r)
			}
		}else{
			writeData(Connection, "\n" + (*response) + "\n")
			log.Println("Session closed by cause: " + (*response))
		}
	}(&response)


	writeData(Connection, WlcMessage + "\n>>")
	for{

		if rcv, err := readData(Connection); rcv != ""{

			if err != nil {
				log.Println(err.Error())
				response = err.Error()

			}else{


				if rcv == "help" || rcv == "\\h"{
					response = h.cli.Help()

				}else if Command, args, err := h.cli.Resolve(rcv); Command != nil{

					response, err = Command.Exec(Tab, args)

					if err != nil{
						response = err.Error()

						// ComQuit
						if err == Com.ErrConnectionClosed{
							return
						}
					}
				}else{
					response = "Unknown command"
					response += "\n" + h.cli.Help()
				}
			}

			writeData(Connection,                                                                                                                                                                                                                                                                                                                                                                                                                                                                       response+ "\n>>")
		}
	}
}



