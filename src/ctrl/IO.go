package ctrl

import (
	"net"
	"strings"
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

	return Connection.Write([]byte(data))
}
