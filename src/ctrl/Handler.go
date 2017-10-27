package ctrl

import (
	"net"
	"fmt"
	"hrentabd"
	"ctrl/controls"
	"io"
)

const WlcMessage = "HrenTab terminal connected OK"

type CommandHandler interface {

	Handle(Connection net.Conn)
}

type commandHandler struct {

	Tab *hrentabd.HrenTab
	t *map[string]string
}

func NewCommandHandler(Tab *hrentabd.HrenTab) CommandHandler{

	ch :=  &commandHandler{Tab:Tab}
	ch.t = &map[string]string{}
	return CommandHandler(ch)
}

func (c *commandHandler)Handle(Connection net.Conn){

	var response string

	defer func(response *string){

		if r := recover(); r != nil{

			if r == io.EOF {
				*response = "client socket closed."
				writeData(Connection, "\n"+(*response)+"\n")
				println("\nSession closed by cause: " + (*response))

			}else{

				writeData(Connection, "\n" + fmt.Sprint(r) + "\n")
				panic(r)
			}
		}else{
			writeData(Connection, "\n" + (*response) + "\n")
			println("\nSession closed by cause: " + (*response))
		}
	}(&response)


	//writeData(Connection, WlcMessage + "\n")
	writeData(Connection, WlcMessage + "\n>>")
	for{

		rcv, err := readData(Connection)

		if err != nil {
			fmt.Println(err.Error())
			response = err.Error()

		}else{

			response = "Unknown command"
			for _,Com := range ComConf{

				if cmdName, args, err := Com.Resolve(rcv); err == nil{

					fmt.Println("Resolved: ", cmdName)

					response, err = Com.Exec(c.Tab, args)

					if err != nil{
						response = err.Error()

						// ComQuit
						if err == controls.ErrConnectionClosed{
							return
						}
					}
					break
				}
			}

		}

		writeData(Connection, ":"+response+ "\n>>")
		//writeData(Connection, "==> " + response + "\n>>")
	}
}


