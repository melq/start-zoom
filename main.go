package main

import (
	"github.com/jessevdk/go-flags"
	"os"
)

type Options struct {
	Start []bool `short:"s" long:"start" description:"Get starting zoom"`
}
var opts Options

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	StartZoomMain(opts)
}
