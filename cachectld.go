package main

import (
	"./cachectl"
	"flag"
	"os"
)

func main() {

	// Parse flags
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		cachectl.PrintVersion(cachectl.Cachectld)
		os.Exit(0)
	}

	// TODO: implement
}
