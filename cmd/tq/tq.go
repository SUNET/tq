package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/SUNET/tq"
)

type Parameters struct {
	url string
}

var p = &Parameters{}

var pubCommand = flag.NewFlagSet("pub", flag.ExitOnError)
var subCommand = flag.NewFlagSet("sub", flag.ExitOnError)

func commonFlags(fs *flag.FlagSet) {
	fs.StringVar(&p.url,"url","tcp://127.0.0.1:9999","endpoint URL")
}

func usage() {
	fmt.Println("usage: tq <command> [<args>]")
        fmt.Println("common arguments:")
        fmt.Println("\t--url=<url>\tspecify the URL to listen/connect to")
        fmt.Println("commands:")
        fmt.Println("\ttq pub [<common arguments>]")
        fmt.Println("\ttq sub [<common arguments>] -- <cmdline>")
        os.Exit(2)
}


func main() {
	tq.Log.Out = os.Stdout
	commonFlags(pubCommand)
	commonFlags(subCommand)

	if len(os.Args) == 1 {
		usage()
	}

	switch os.Args[1] {
	case "pub":
		pubCommand.Parse(os.Args[2:])
		tq.Sub(p.url, flag.Args())
	case "sub":
		subCommand.Parse(os.Args[2:])
		tq.Pub(p.url)
	case "help",
	     "--help",
	     "?":
		usage()
	default:
		fmt.Println("%q is not a valid command", os.Args[1])
		os.Exit(2)
	}
}
