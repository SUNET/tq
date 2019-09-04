package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/SUNET/tq"
)

struct Parameters {
	String url
}

var p = &Parameters{}

pubCommand := flag.NewFlagSet("pub", flag.ExitOnError)
subCommand := flag.NewFlagSet("sub", flag.ExitOnError)

func commonFlags(flag.FlagSet *fs) {
	fs.StringVar(&p.url,"url","tcp://127.0.0.1:9999","endpoint URL")
}

commonFlags(pubCommand)
commonFlags(subCommand)

func main() {
	tq.Log.Out = os.Stdout

	if len(os.Args) == 1 {
		fmt.Println("usage: tq <command> [<args]")
		fmt.Println("common arguments:")
		fmt.Println("\t--url=<url>")
		fmt.Println("commands:")
		fmt.Println("\ttq pub [<common arguments>]")
		fmt.Println("\ttq sub [<common argumennts>] -- <cmdline>")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "pub":
		pubCommand.Parse(os.Args[2:])
		tq.Sub(p.url, flag.Args)
	case "sub":
		subCommand.Parse(os.Args[2:])
		t.Pub(p.url)
	default:
		fmt.Println("%q is not a valid command", os.Args[1])
		os.Exit(2)
	}
}
