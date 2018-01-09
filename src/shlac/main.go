package main

import (
	"os"
	"github.com/urfave/cli"
	"fmt"
	"errors"
	"shared/sig"
	"net"
	"bufio"
	"regexp"
	"strings"
	"io/ioutil"
	"encoding/json"
	"bytes"
	. "shared/config"
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
	app.Version             = "0.1"
	app.Name                = "SHLAC"
	app.Usage               = "SHlac Like As Cron"
	app.Author              = "Evgeny Nefedkin"
	app.Email               = "evgeny.nefedkin@umbrella-web.com"
	app.EnableBashCompletion= true
	app.Description         = "Distributed and concurrency job manager\n" +

		"\t\tSupported extended syntax:\n" +
		"\t\t------------------------------------------------------------------------\n" +
		"\t\tField name     Mandatory?   Allowed values    Allowed special characters\n" +
		"\t\t----------     ----------   --------------    --------------------------\n" +
		"\t\tSeconds        No           0-59              * / , -\n" +
		"\t\tMinutes        Yes          0-59              * / , -\n" +
		"\t\tHours          Yes          0-23              * / , -\n" +
		"\t\tDay of month   Yes          1-31              * / , - L W\n" +
		"\t\tMonth          Yes          1-12 or JAN-DEC   * / , -\n" +
		"\t\tDay of week    Yes          0-6 or SUN-SAT    * / , - L #\n" +
		"\t\tYear           No           1970–2099         * / , -\n" +

		"\n\n" +

		"\t\tand aliases:\n" +
		"\t\t-------------------------------------------------------------------------------------------------\n" +
		"\t\tEntry       Description                                                             Equivalent to\n" +
		"\t\t-------------------------------------------------------------------------------------------------\n" +
		"\t\t@annually   Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *\n" +
		"\t\t@yearly     Run once a year at midnight in the morning of January 1                 0 0 0 1 1 * *\n" +
		"\t\t@monthly    Run once a month at midnight in the morning of the first of the month   0 0 0 1 * * *\n" +
		"\t\t@weekly     Run once a week at midnight in the morning of Sunday                    0 0 0 * * 0 *\n" +
		"\t\t@daily      Run once a day at midnight                                              0 0 0 * * * *\n" +
		"\t\t@hourly     Run once an hour at the beginning of the hour                           0 0 * * * * *\n" +
		"\t\t@reboot     Not supported"


	// CONFIG
	app.Flags =  []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: `/etc/shlacd/shlacd.conf`,
			Usage: "path to daemon config-file",
		},
	}


	// COMMANDS
	app.Commands = []cli.Command{
		{// IMPORT
			Name:    "import",
			Aliases: []string{"i"},
			Usage:   "import jobs from cron-formatted file",
			UsageText: "Example: " +
				"shlac import <path/to/import/file>",

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
				if c.Bool("purge"){ purge(connection) }

				Import(filePath, connection, !c.Bool("skip-check"))

				return nil
			},
		},

		{// ADD JOB
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add job from cron-formatted line",
			UsageText: "Example: " +
				"shlac add '<cron-formatted line>'",

			Flags: 	[]cli.Flag{
				cli.BoolFlag{
					Name:  "skip-check, s",
					Usage: "add job even if same is already exist (skip checking for duplicates)",
				},
			},

			Action:  func(c *cli.Context) error {

				cronString := c.Args().Get(0)
				if cronString == "" {
					panic(ErrCmdArgs)
				}

				connection := connect( loadConfig(c.GlobalString("config")) )
				defer func(){
					connection.Write([]byte(`\q`))
					connection.Close()
				}()


				ImportLine(cronString, connection, !c.Bool("skip-check"))

				return nil
			},
		},

		{// EXPORT
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


				Export(filePath, connection)

				return nil
			},
		},
	}

	app.Run(os.Args)
}


func Export(filePath string, connection net.Conn){

	flushConnection(connection) // clear socket buffer
	connection.Write([]byte(`\l`))

	response := flushConnection(connection) // read answer

	response = response[:len(response)-4] // remove terminal bytes

	re := regexp.MustCompile(`(?m)^.+?\s+`) // remove jobs id
	response = re.ReplaceAll(response, []byte{})

	fmt.Println(string(response))


	ioutil.WriteFile(filePath, response, 0644)
}

func Import(filePath string, connection net.Conn, checkDuplicates bool){

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner     := bufio.NewScanner(file)

	for scanner.Scan() {

		ImportLine(scanner.Text(), connection, checkDuplicates)

		// MAIN LOOP
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}

func ImportLine(cronString string, connection net.Conn, checkDuplicates bool){

	delimiter   := regexp.MustCompile(`\s+`)

	parts := delimiter.Split(cronString, 6)

	cronLine    := strings.Join(parts[:5], " ")
	commandLine := parts[5]

	importLine := fmt.Sprintf(`\a -cron "%s" -cmd "%s"`+"\n", cronLine, commandLine)

	if cronLine[:1] == `#`{
		fmt.Printf("SKIPP (disabled)>> %s", importLine)
		return
	}

	if checkDuplicates && isDuplicated(cronString, connection){
		fmt.Printf("SKIPP (duplicated)>> %s", importLine)
		return
	}


	fmt.Printf("IMPORT>> %s", importLine)
	connection.Write([]byte(importLine))
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

func flushConnection(connection net.Conn) (flushed []byte){

	bufSize := 256
	buf := make([]byte, bufSize)

	for{
		n,e := connection.Read(buf)

		flushed = append(flushed, buf[:n]...)

		if e != nil || n < bufSize {break}
	}

	return flushed
}

func purge(connection net.Conn){
	connection.Write([]byte(`\r --all`))
}

func isDuplicated(cronLine string, connection net.Conn) bool {

	cronLine = strings.Replace(cronLine, `"`, `\"`, -1)

	flushConnection(connection)
	connection.Write([]byte(`\g -c "` +cronLine+ `"`))

	response := make([]byte, 8)
	connection.Read(response)

	return !bytes.Equal(response, []byte{110, 117, 108, 108, 0, 10, 62, 62})
}