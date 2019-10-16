# static-server

[![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](LICENSE)  
[![Go Report Card](https://goreportcard.com/badge/github.com/dictyBase/static-server)](https://goreportcard.com/report/github.com/dictyBase/static-server)
[![Technical debt](https://badgen.net/codeclimate/tech-debt/dictyBase/static-server)](https://codeclimate.com/github/dictyBase/static-server/trends/technical_debt)
[![Issues](https://badgen.net/codeclimate/issues/dictyBase/static-server)](https://codeclimate.com/github/dictyBase/static-server/issues)
[![Maintainability percentage](https://badgen.net/codeclimate/maintainability-percentage/dictyBase/static-server)](https://codeclimate.com/github/dictyBase/static-server)
[![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=dictyBase/static-server)](https://dependabot.com)  
![GitHub repo size](https://img.shields.io/github/repo-size/dictyBase/static-server?style=plastic)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/dictyBase/static-server?style=plastic)
[![Lines of Code](https://badgen.net/codeclimate/loc/dictyBase/static-server)](https://codeclimate.com/github/dictyBase/static-server/code)  
![Commits](https://badgen.net/github/commits/dictyBase/static-server/develop)
![Last commit](https://badgen.net/github/last-commit/dictyBase/static-server/develop)
![Branches](https://badgen.net/github/branches/dictyBase/static-server)
![Tags](https://badgen.net/github/tags/dictyBase/static-server)  
![Issues](https://badgen.net/github/issues/dictyBase/static-server)
![Open Issues](https://badgen.net/github/open-issues/dictyBase/static-server)
![Closed Issues](https://badgen.net/github/closed-issues/dictyBase/static-server)
![Total PRS](https://badgen.net/github/prs/dictyBase/static-server)
![Open PRS](https://badgen.net/github/open-prs/dictyBase/static-server)
![Closed PRS](https://badgen.net/github/closed-prs/dictyBase/static-server)
![Merged PRS](https://badgen.net/github/merged-prs/dictyBase/static-server)  
[![Funding](https://badgen.net/badge/NIGMS/Rex%20L%20Chisholm,dictyBase/yellow?list=|)](https://projectreporter.nih.gov/project_info_description.cfm?aid=9476993)
[![Funding](https://badgen.net/badge/NIGMS/Rex%20L%20Chisholm,DSC/yellow?list=|)](https://projectreporter.nih.gov/project_info_description.cfm?aid=9438930)

A golang based file server with default logrus based logging.

# Available commands

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

## Developers

<a href="https://sourcerer.io/cybersiddhu"><img src="https://sourcerer.io/assets/avatar/cybersiddhu" height="80px" alt="Sourcerer"></a>
<a href="https://sourcerer.io/wildlifehexagon"><img src="https://sourcerer.io/assets/avatar/wildlifehexagon" height="80px" alt="Sourcerer"></a>
