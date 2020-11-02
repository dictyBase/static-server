package main

import (
	"log"
	"os"

	"github.com/dictyBase/static-server/commands"
	"github.com/dictyBase/static-server/validate"
	"github.com/urfave/cli"
)

var staticF = `The static files will only be served from this static folder
		  and expected to be under the base folder. The url path should
		  also match the filesystem. Any other path will
		  be redirected to the index.html
`

func main() {
	app := cli.NewApp()
	app.Name = "static-server"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:   "run",
			Usage:  "A http static file server",
			Action: commands.ServeAction,
			Before: validate.ValidateServer,
			Flags:  serverFlags(),
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func serverFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "folder, f",
			Usage:  "Location of folder from where files will be served[required]",
			EnvVar: "FILE_FOLDER",
		},
		cli.IntFlag{
			Name:  "port, p",
			Usage: "http port, default is 9595",
			Value: 9595,
		},
		cli.StringFlag{
			Name:   "log-format",
			Usage:  "log format, json or text",
			EnvVar: "LOG_FORMAT",
			Value:  "json",
		},
		cli.StringFlag{
			Name:   "log-file, l",
			Usage:  "Name of the log file, default goes to stderr",
			EnvVar: "LOG_FILE",
		},
		cli.StringFlag{
			Name:   "sub-url",
			Usage:  "Alternate url path that does not match the filesystem",
			EnvVar: "SUB_URL",
		},
		cli.StringFlag{
			Name:   "static-folder,sf",
			Usage:  staticF,
			EnvVar: "STATIC_FOLDER",
			Value:  "/static",
		},
		cli.IntFlag{
			Name:  "cache-duration,d",
			Usage: "how long the static assets will be cached given in months",
			Value: 11,
		},
	}
}
