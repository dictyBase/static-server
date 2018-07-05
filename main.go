package main

import (
	"os"

	"github.com/dictyBase/static-server/commands"
	"github.com/dictyBase/static-server/validate"
	"gopkg.in/urfave/cli.v1"
)

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
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "folder, f",
					Usage:  "Location of folder from where files will be served[required]",
					EnvVar: "STATIC_FOLDER",
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
					Usage:  "Alternate url path if the files are not being served from root",
					EnvVar: "SUB_URL",
				},
				cli.StringFlag{
					Name:   "virtual-static-folder,vf",
					Usage:  "virtual static folder if it does not match the physical path",
					EnvVar: "VIRTUAL_STATIC_FOLDER",
				},
			},
		},
	}
	app.Run(os.Args)
}
