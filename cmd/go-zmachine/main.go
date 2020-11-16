package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

type options struct {
	version int
}

func newOptions() *options {
	return &options{}
}

func (o *options) Configure() {
	flag.IntVar(&o.version, "version", o.version, "Force a specific Z-machine version (1-8).")

}

func (o *options) Validate() error {
	if o.visited("version") && (o.version < 1 || o.version > 8) {
		return fmt.Errorf("%d is not a valid Z-machine version", o.version)
	}

	return nil
}

func (o *options) Version() (version int, wasSet bool) {
	return o.version, o.visited("version")
}

func (o *options) visited(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func usage() {
	fmt.Print(`
Welcome to the goz Z-machine interpreter!

Usage:
  goz <opts> <filename>

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage

	opts := newOptions()
	opts.Configure()

	flag.Parse()
	if err := opts.Validate(); err != nil {
		fmt.Printf("error: %+v\n", err)
		return
	}

	if len(flag.Args()) != 1 {
		fmt.Printf("error: wrong number of arguments\n")
		return
	}

	gamefile := flag.Arg(0)
	fileContents, err := ioutil.ReadFile(gamefile)
	if err != nil {
		fmt.Printf("read game file: %+v", err)
	}

	var gamefileVersion int

	if forcedVersion, wasSet := opts.Version(); wasSet {
		gamefileVersion = forcedVersion
	} else {
		gamefileVersion, err = strconv.Atoi(gamefile[len(gamefile)-1:])
		if err != nil {
			fmt.Printf("extract game version: %+v", err)
		}
	}

	fmt.Printf("%s is a z%d gamefile weighing in at %d bytes", gamefile, gamefileVersion, len(fileContents))
}
