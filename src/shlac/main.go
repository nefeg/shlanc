package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"errors"
	"sig"
	"net"
	"bufio"
	"regexp"
	"strings"
	"io/ioutil"
	"encoding/json"
	"bytes"
	. "config"
)

var ErrCmdArgs      = errors.New("ERR: expected argument")
var ErrConfFile     = errors.New("ERR: invalid config file")
var ErrConfInvalid  = errors.New("ERR: invalid config")

func init()  {
	sig.SIG_INT(nil)
}


func main(){

	defer func(){
		if r := recover(); r != nil{

			fmt.Println(r)

			if r == ErrCmdArgs{
				fmt.Println("See: shlac <command> --help")
			}
		}
	}()



	app := cli.NewApp()
	app.Version     = "0.1"
	app.Name        = "ShLAC"
	app.Usage       = "SHlac Like As Cron"
	app.Author      = "Evgeny Nefedkin"
	app.Email       = "evgeny.nefedkin@umbrella-web.com"
	app.Description = "Distributed and concurrency job manager"

	// CONFIG
	app.Flags =  []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: `/etc/shlac/shlac.conf`,
			Usage: "path to daemon config-file",
		},
	}


	// COMMANDS
	app.Commands = []cli.Command{
		{
			Name:    "import",
			Aliases: []string{"i"},
			Usage:   "import jobs from cron-formatted file",
			UsageText: "Example: shlac import <path/to/import/file>",

			Flags: 	[]cli.Flag{
				cli.BoolFlag{
					Name:  "purge, p",
					Usage: "delete jobs before import",
				},

				cli.BoolFlag{
					Name:  "skip-check, s",
					Usage: "add job even if same is already exist (skip checking for duplicates)",
				},
			},



			Action:  func(c *cli.Context) error {

				filePath := c.Args().Get(0)
				if filePath == "" {
					panic(ErrCmdArgs)
				}

				connection := connect( loadConfig(c.GlobalString("config")) )
				defer func(){
					connection.Write([]byte(`\q`))
					connection.Close()
				}()

				// clean table before import
				if c.Bool("purge"){ _purge(connection) }

				_import(filePath, connection, !c.Bool("skip-check"))

				return nil
			},
		},
		{
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "export jobs to file in cron-format",
			UsageText: "Example: shlac export <path/to/export/file>",
			Action:  func(c *cli.Context) error {

				filePath := c.Args().Get(0)
				if filePath == "" {
					panic(ErrCmdArgs)
				}

				connection := connect( loadConfig(c.GlobalString("config")) )
				defer func(){
					connection.Write([]byte(`\q`))
					connection.Close()
				}()


				_export(filePath, connection)

				return nil
			},
		},
	}

	app.Run(os.Args)
}


func loadConfig(configPath string) (config *Config) {

	configRaw, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(ErrConfFile)
	}

	config = &Config{}
	if err := json.Unmarshal(configRaw, config); err != nil{
		panic(ErrConfInvalid)
	}

	return config
}

func connect(config *Config) (connection net.Conn){
	if config.Client.Type != "socket" {
		panic("Unsupported client type")
	}

	conn, err := net.Dial(config.Client.Options.Network, config.Client.Options.Address)
	if err != nil{
		panic(err)
	}

	return conn
}

func clearConnection(connection net.Conn) (l int){

	buf := make([]byte, 4096) // big buffer
	n, _ := connection.Read(buf)

	return n
}


func _export(filePath string, connection net.Conn){

	clearConnection(connection)
	connection.Write([]byte(`\l`))

	var response []byte
	bufSize := 256
	buf := make([]byte, bufSize)

	for{
		n,e := connection.Read(buf)

		response = append(response, buf[:n]...)

		if e != nil || n < bufSize {break}
	}

	response = response[:len(response)-4] // remove terminal bytes

	re := regexp.MustCompile(`(?m)^.+?\s+`) // remove jobs id
	response = re.ReplaceAll(response, []byte{})

	fmt.Println(string(response))


	ioutil.WriteFile(filePath, response, 0644)
}

func _import(filePath string, connection net.Conn, checkDuplicates bool){

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner     := bufio.NewScanner(file)
	delimiter   := regexp.MustCompile(`\s+`)

	for scanner.Scan() {
		parts := delimiter.Split(scanner.Text(), 6)

		cronLine    := strings.Join(parts[:5], " ")
		commandLine := parts[5]

		importLine := fmt.Sprintf(`\a -cron "%s" -cmd "%s"`+"\n", cronLine, commandLine)

		if cronLine[:1] == `#`{
			fmt.Printf("SKIPP (disabled)>> %s", importLine)
			continue
		}

		if checkDuplicates && __isDuplicated(scanner.Text(), connection){
			fmt.Printf("SKIPP (duplicated)>> %s", importLine)
			continue
		}


		fmt.Printf("IMPORT>> %s", importLine)
		connection.Write([]byte(importLine))

		// MAIN LOOP
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}


func _purge(connection net.Conn){
	connection.Write([]byte(`\r --all`))
}

func __isDuplicated(cronLine string, connection net.Conn) bool {

	cronLine = strings.Replace(cronLine, `"`, `\"`, -1)

	clearConnection(connection)
	connection.Write([]byte(`\g -c "` +cronLine+ `"`))

	response := make([]byte, 8)
	connection.Read(response)

	return !bytes.Equal(response, []byte{110, 117, 108, 108, 0, 10, 62, 62})
}