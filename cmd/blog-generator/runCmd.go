package main

import (
	"flag"
	"net/http"
	"strconv"
)

const runShortHelp = `Runs a web server for the generated blog`
const runLongHelp = `
Runs a web server for the generated blog.

The default port for the server is 8080.

`

type runCmd struct {
	source string
	port   int
}

func (cmd *runCmd) Name() string      { return "run" }
func (cmd *runCmd) Args() string      { return "" }
func (cmd *runCmd) ShortHelp() string { return runShortHelp }
func (cmd *runCmd) LongHelp() string  { return runLongHelp }
func (cmd *runCmd) Hidden() bool      { return false }

func (cmd *runCmd) Register(fs *flag.FlagSet) {
	fs.IntVar(&cmd.port, "p", 8080, "port for web server")
}

func (cmd *runCmd) Run(ctx *Ctx, args []string) error {
	fs := http.FileServer(http.Dir(cmd.source))
	http.Handle("/", fs)

	serverPort := strconv.Itoa(cmd.port)
	ctx.Out.Println("Listening @ localhost:" + serverPort + "/")
	http.ListenAndServe(":"+serverPort, nil)
	return nil
}
