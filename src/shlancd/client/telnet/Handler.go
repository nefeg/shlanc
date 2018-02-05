package telnet

import (
	"net"
	"os"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"github.com/urfave/cli"

	capi "shlancd/cli"
	"github.com/mattn/go-shellwords"
	"io"
	"regexp"
	"errors"
)

type handler struct {
	addr    net.Addr
}

const logPrefix = "[client.telnet] "
const WlcMessage =  "HrenTab terminal connected OK\n" +
					"Type \"help\" or \"\\h\" for show available commands"

var ErrConnectionClosed = errors.New("** command <QUIT> received")


func New(listen net.Addr) *handler{

	return &handler{ addr:listen}
}


func (h *handler) Handle(context capi.Context){

	IPC, err := net.Listen(h.addr.Network(), h.addr.String())
	if err != nil {
		slog.PanicF("%s %s %s",logPrefix,  "panic:", err.Error())
	}
	slog.InfoF(logPrefix + "Listening: %s://%s\n", IPC.Addr().Network(), IPC.Addr().String())
	slog.InfoLn(logPrefix + "Connection ID: ", &IPC)


	defer func(){
		IPC.Close()
		if UAddr, err := net.ResolveUnixAddr(h.addr.Network(), h.addr.String()); err == nil{
			os.Remove(UAddr.String())
		}
	}()

	for{
		if Connection, err := IPC.Accept(); err == nil {

			go func(){
				slog.InfoF(logPrefix + "New client connection accepted [connection ID: %v]\n", &Connection)

				h.handleConnection(Connection, context)
				Connection.Close()

				slog.InfoF(logPrefix + "Client connection closed [connection ID: %v]\n", &Connection)
			}()

		}else{
			slog.CritLn(err.Error())
			continue
		}
	}
}

func (h *handler)handleConnection(Connection net.Conn, context capi.Context){

	var response string
	var Responder = NewClient(Connection)

	defer func(response *string){

		if r := recover(); r != nil{

			if r == io.EOF {
				*response = "client socket closed."
				Responder.WriteString("\n" + (*response) + "\n")
				slog.InfoLn(logPrefix + "Session closed by cause: " + (*response))

			}else{
				slog.InfoLn(logPrefix + "Session closed by cause: " , r)
			}
		}else{
			Responder.WriteString("\n" + (*response) + "\n")
			slog.InfoLn(logPrefix + "Session closed by cause: " + (*response))
		}
	}(&response)


	Responder.WriteString(WlcMessage)


	Cli := capi.New()

	Cli.Writer              = Responder
	Cli.ErrWriter           = Responder

	// COMMANDS
	Cli.Commands = []cli.Command{
		capi.NewComAdd(&context),
		capi.NewComList(&context),
		capi.NewComRemove(&context),
		capi.NewComPurge(&context),
		capi.NewComGet(&context),
		{
			Name:    "exit",
			Aliases: []string{`q`},
			Usage:   "close connection",
			UsageText: "Example: " ,

			Action:  func(c *cli.Context) error {

				slog.DebugLn("Action: exit")

				c.App.Writer.Write([]byte("Sending <QUIT> signal..."))
				panic(ErrConnectionClosed)

				return nil
			},
		},
	}

	Cli.After = func(c *cli.Context) error {

		c.App.Writer.Write(PacketTerm)

		return nil
	}

	Cli.ExitErrHandler = func(c *cli.Context, err error){
		c.App.Writer.Write([]byte(err.Error()))
		slog.DebugLn(logPrefix, err)
	}


	for{

		if rcb, err := Responder.ReadData(); len(rcb) != 0{

			if err != nil {
				slog.CritLn(err.Error())
				response = err.Error()
				Responder.WriteString(response)

			}else{

				slog.DebugLn(logPrefix, "Args (byte,raw):", rcb)
				slog.DebugLn(logPrefix, "Args (string,raw):", string(rcb))

				if match,_ := regexp.Match(`^\w.*`, rcb); match != true{
					rcb = []byte("help")
					slog.DebugLn(logPrefix, "Incorrect args, show help")
				}

				args,_ := shellwords.Parse("self " + string(rcb))

				Cli.Run( args )

				slog.DebugLn(logPrefix, "Cli.Run (complete)")
			}
		}
	}
}



