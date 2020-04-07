package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/SUNET/tq/pkg/pipeline"
	"github.com/sirupsen/logrus"
	"github.com/spy16/sabre"
	"github.com/spy16/sabre/repl"
)

var Log = logrus.New()

var helpFlag bool
var branch string
var commit string
var version string

func usage(code int) {
	fmt.Println("usage: tq [-h] [-e <expression>]")
	os.Exit(code)
}

func ver() string {
	if len(branch) > 0 && len(commit) > 0 {
		return fmt.Sprintf("%s@%s", commit, branch)
	} else if len(version) > 0 {
		return fmt.Sprintf("v%s", version)
	} else {
		return "unknown"
	}
}

func main() {
	Log.Out = os.Stdout
	flag.Parse()
	if helpFlag {
		usage(0)
	}
	scope := sabre.NewScope(nil)
	for k, v := range pipeline.PipelineFactories {
		_ = scope.BindGo(k, v)
	}
	stat, _ := os.Stdin.Stat()
	defer func() {
		if r := recover(); r != nil {
			Log.Debug("recovered...")
		}
	}()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println("data is being piped to stdin")
	} else {
		repl.New(scope,
			repl.WithBanner(fmt.Sprintf("tq shell [%s]", ver())),
			repl.WithPrompts(">", "|"),
		).Loop(context.Background())
	}
}

func init() {
	flag.BoolVar(&helpFlag, "h", false, "show help")
}
