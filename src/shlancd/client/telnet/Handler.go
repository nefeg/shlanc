package telnet

import (
	"net"
	"os"
	"github.com/umbrella-evgeny-nefedkin/slog"
	"github.com/urfave/cli"

	capi "shlancd/cli"
	"github.com/mattn/go-shellwords"
	"io"
	"errors"
	"regexp"
)

type handler struct {
	addr    net.Addr
}

const logPrefix = "[client.telnet] "
const WlcMessage =  "ShlaNc terminal connected OK\n" +
					"server version: 0.23; client version: 0.3\n" +
					"Type \"help\" for show available commands"

var ErrConnectionClosed = errors.New("** command <QUIT> received")


func New(listen net.Addr) *handler{

	return &handler{ addr:listen}
}


func (h *handler) Handle(context capi.Context){

	IPC, err := net.Listen(h.addr.Network(), h.addr.String())
	if err != nil {
		slog.Fatalln(logPrefix + " panic: ", err.Error())
	}
	slog.Infof(logPrefix + "Listening: %s://%s\n", IPC.Addr().Network(), IPC.Addr().String())
	slog.Infoln(logPrefix + "Connection ID: ", &IPC)


	defer func(){
		IPC.Close()
		if UAddr, err := net.ResolveUnixAddr(h.addr.Network(), h.addr.String()); err == nil{
			os.Remove(UAddr.String())
		}
	}()

	for{
		if Connection, err := IPC.Accept(); err == nil {

			go func(){
				slog.Infof(logPrefix + "New client connection accepted [connection ID: %v]\n", &Connection)

				h.handleConnection(Connection, context)
				Connection.Close()

				slog.Infof(logPrefix + "Client connection closed [connection ID: %v]\n", &Connection)
			}()

		}else{
			slog.Critln(err.Error())
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
				slog.Infoln(logPrefix + "Session closed by cause: " + (*response))

			}else{
				slog.Infoln(logPrefix + "Session closed by cause: " , r)
			}
		}else{
			Responder.WriteString("\n" + (*response) + "\n")
			slog.Infoln(logPrefix + "Session closed by cause: " + (*response))
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

				slog.Debugln("Action: exit")

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
		slog.Debugln(logPrefix, err)
	}


	for{

		if rcb, err := Responder.ReadData(); len(rcb) != 0{

			if err != nil {
				slog.Critln(err.Error())
				response = err.Error()
				Responder.WriteString(response)

			}else{

				slog.Debugln(logPrefix, "Args (byte,raw):", rcb)
				slog.Debugln(logPrefix, "Args (string,raw):", string(rcb))

				if match,_ := regexp.Match(`^\w.*`, rcb); match != true{
					rcb = []byte("help")
					slog.Debugln(logPrefix, "Incorrect args, show help")
				}

				args,_ := shellwords.Parse("self " + string(rcb))

				Cli.Run( args )

				slog.Debugln(logPrefix, "Cli.Run (complete)")
			}
		}
	}
}



