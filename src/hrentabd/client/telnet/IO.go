package telnet

import (
	"net"
	"strings"
	"log"
)

func readData(Connection net.Conn) (rcv string, err error){

	tmp := make([]byte, 4096)

	length, err := Connection.Read(tmp[:])
	if err != nil {
		panic(err)
	}


	if length>0{
		rcv = strings.TrimSpace(string(tmp[:length]))
	}

	return rcv, err
}

func writeData(Connection net.Conn, data string) (int, error){

	response := append([]byte(data), []byte{00, 10, 62, 62}...)

	log.Println("[SYS]writeData:", response)

	return Connection.Write(response)
}
