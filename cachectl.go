package main

import (
	"./cachectl"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
)

func main() {

	// Parse flags
	version := flag.Bool("v", false, "show version")
	op := flag.String("op", "stat", "operation(stat, purge)")
	fpath := flag.String("f", "", "target file path")
	filter := flag.String("filter", "*", "filter pattern")
	rate := flag.Float64("r", 1.0, "rate of page cache purged(0.0 <= r <= 1.0)")
	verbose := flag.Bool("verbose", false, "verbose mode")
	flag.Parse()

	if *version {
		cachectl.PrintVersion(cachectl.Cachectl)
		os.Exit(0)
	}

	if *fpath == "" {
		log.Println("target file path is empty.")
		os.Exit(0)
	}

	fi, err := os.Stat(*fpath)
	if err != nil {
		log.Fatal(err.Error())
	}

	if *filter == "*" {
		*filter = ".*"
	}

	re := regexp.MustCompile(*filter)

	if *op == "stat" {
		if fi.IsDir() {
			err := cachectl.WalkPrintPagesStat(*fpath, re)
			if err != nil {
				log.Fatal(fmt.Sprintf("failed to walk in %s.", fi.Name()))
			}
		} else {
			if !fi.Mode().IsRegular() {
				log.Fatal(fmt.Sprintf("%s is not regular file", fi.Name()))
			}

			cachectl.PrintPagesStat(*fpath, fi.Size())
		}
	} else {
		if fi.IsDir() {
			err := cachectl.WalkPurgePages(*fpath, re, *rate, *verbose)
			if err != nil {
				log.Fatal(fmt.Sprintf("failed to walk in %s.", fi.Name()))
			}
		} else {
			if !fi.Mode().IsRegular() {
				log.Fatal(fmt.Sprintf("%s is not regular file", fi.Name()))
			}

			cachectl.RunPurgePages(*fpath, fi.Size(), *rate, *verbose)
		}
	}
}
