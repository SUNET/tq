package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spy16/sabre"
	"github.com/spy16/sabre/repl"
	"github.com/sunet/tq/pkg/api"
	"github.com/sunet/tq/pkg/message"
	"github.com/sunet/tq/pkg/meta"
	"github.com/sunet/tq/pkg/pipeline"
)

var Log = logrus.New()

var helpFlag bool
var relpFlag bool
var logLevelFlag string

func usage(code int) {
	fmt.Println("usage: tq [-h] [-e <expression>]")
	os.Exit(code)
}

func is_not_tty() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func configLogger(log *logrus.Logger, ll string) {
	log.Out = os.Stdout

	if len(ll) > 0 {
		level, err := logrus.ParseLevel(logLevelFlag)
		if err != nil {
			log.Panicf("Unable to parse loglevel: %s", err.Error())
		}
		log.SetLevel(level)
	}
}

func readEvalFiles(scope sabre.Scope, files ...string) {
	srf := NewScriptReaderFactory()
	for _, g := range files {
		matches, _ := filepath.Glob(g)
		for _, r := range matches {
			Log.Debugf("About to load %s", r)
			f, err := os.Open(r)
			defer f.Close()
			if err != nil {
				Log.Fatalf("Unable to open %s: %s", r, err.Error())
			}
			_, err = srf.ReadEval(scope, bufio.NewReader(f))
			if err != nil {
				Log.Fatalf("Unable to execute %s: %s", r, err.Error())
			}
		}
	}
}

func main() {

	flag.Parse()
	if helpFlag {
		usage(0)
	}

	configLogger(Log, logLevelFlag)
	configLogger(message.Log, logLevelFlag)
	configLogger(pipeline.Log, logLevelFlag)
	configLogger(api.Log, logLevelFlag)

	defer func() {
		if r := recover(); r != nil {
			Log.Debug(r)
		}
	}()

	files := flag.Args()
	relpFlag = relpFlag || (len(files) == 0)

	scope := SabreScope()
	srf := NewScriptReaderFactory()

	if relpFlag {
		repl.New(scope,
			repl.WithBanner(fmt.Sprintf("tq shell [%s]", meta.Version())),
			repl.WithPrompts(">", "|"),
			repl.WithReaderFactory(srf),
		).Loop(context.Background())
	} else {
		readEvalFiles(scope, files...)

		if is_not_tty() {
			_, err := srf.ReadEval(scope, os.Stdin)
			if err != nil {
				Log.Fatalf("Unable to execute from stdin: %s", err.Error())
			}
		}
	}
}

func init() {
	flag.BoolVar(&helpFlag, "h", false, "show help")
	flag.BoolVar(&relpFlag, "s", false, "execute RELP (read-eval-print) loop")
	flag.StringVar(&logLevelFlag, "loglevel", "info", "loglevel")
}
