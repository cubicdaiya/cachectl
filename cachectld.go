package main

import (
	"./cachectl"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	// Parse flags
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		cachectl.PrintCachectlVersion()
		os.Exit(0)
	}

	// TODO: implement
}
