package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/spy16/slurp"
	"github.com/spy16/slurp/core"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spy16/slurp/repl"
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

func isNotTty() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func ConfigLoggers(logLevelFlag string) {
	configLogger(Log, logLevelFlag)
	configLogger(message.Log, logLevelFlag)
	configLogger(pipeline.Log, logLevelFlag)
	configLogger(api.Log, logLevelFlag)
}

func configLogger(log *logrus.Logger, ll string) {
	log.SetOutput(os.Stdout)

	if len(ll) > 0 {
		level, err := logrus.ParseLevel(logLevelFlag)
		if err != nil {
			log.Panicf("Unable to parse loglevel: %s", err.Error())
		}
		log.SetLevel(level)
	}
}

func ReadEval(sl *slurp.Interpreter, r io.Reader) (core.Any, error) {
	var data []byte
	mod, err := r.Read(data)
	if err != nil {
		return nil, err
	}

	return sl.Eval(mod)
}

func readEvalFiles(sl *slurp.Interpreter, files ...string) core.Any {
	var v core.Any
	for _, g := range files {
		matches, _ := filepath.Glob(g)
		for _, r := range matches {
			Log.Debugf("About to load %s", r)

			s, err := ioutil.ReadFile(r)
			v, err = sl.EvalStr(string(s))
			if err != nil {
				Log.Fatalf("Unable to execute %s: %s", r, err.Error())
			}
		}
	}
	return v
}

func main() {

	flag.Parse()
	if helpFlag {
		usage(0)
	}

	ConfigLoggers(logLevelFlag)

	defer func() {
		if r := recover(); r != nil {
			Log.Debug(r)
		}
	}()

	files := flag.Args()
	relpFlag = relpFlag || (len(files) == 0)
	srf := NewScriptReaderFactory()
	sl := NewInterpreter()

	if relpFlag {
		repl.New(sl,
			repl.WithBanner(fmt.Sprintf("tq shell [%s]", meta.Version())),
			repl.WithPrompts(">", "|"),
			repl.WithReaderFactory(srf),
		).Loop(context.Background())
	} else {
		readEvalFiles(sl, files...)

		if isNotTty() {
			_, err := srf.ReadEval(sl, os.Stdin)
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
