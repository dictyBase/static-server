# static-server
A golang based file server with default logrus based logging.

# Availabe commands 

```
NAME:
   static-server - A new cli application

USAGE:
   static-server [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     run      A http static file server 
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Subcommands

```
NAME:
   static-server run - A http static file server

USAGE:
   static-server run [command options] [arguments...]

OPTIONS:
   --folder value, -f value           Location of folder from where files will be served[required] [$FILE_FOLDER]
   --port value, -p value             http port, default is 9595 (default: 9595)
   --log-format value                 log format, json or text (default: "json") [$LOG_FORMAT]
   --log-file value, -l value         Name of the log file, default goes to stderr [$LOG_FILE]
   --sub-url value                    Alternate url path that does not match the filesystem [$SUB_URL]
   --static-folder value, --sf value  The static files will only be served from this static folder
                                          and expected to be under the base folder. The url path should
                                          also match the filesystem. Any other path will
                                          be redirected to the index.html(default: "/static") [$STATIC_FOLDER]
```
