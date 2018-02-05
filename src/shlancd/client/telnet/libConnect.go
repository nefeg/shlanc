package telnet

import (
	"net"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"bytes"
)

var PacketTerm = []byte{00, 10, 62, 62}

func readData(connection net.Conn) (data []byte, err error){

	tmp := make([]byte, 4096)

	length, err := connection.Read(tmp[:])
	if err != nil {
		slog.PanicLn(err)
	}

	slog.DebugLn("[client.telnet] readData (byte,raw):", tmp[:length])

	if length>0{
		data = bytes.TrimSpace(tmp[:length])
	}

	slog.DebugLn("[client.telnet] readData (byte,trim):", tmp[:length])
	slog.DebugLn("[client.telnet] readData (string,trim):", string(data))

	return data, err
}

func writeData(connection net.Conn, data []byte) (int, error){

	slog.DebugLn("[client.telnet] writeData (raw string):", data)

	response := append(data, PacketTerm...)

	slog.DebugLn("[client.telnet] writeData (bytes):", response)

	return connection.Write(response)
}
