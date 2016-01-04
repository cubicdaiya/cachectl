package main

import (
	"flag"
	"github.com/cubicdaiya/cachectl/cachectl"
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
		return
	}

	if *fpath == "" {
		flag.Usage()
		return
	}

	fi, err := os.Stat(*fpath)
	if err != nil {
		log.Fatal(err)
	}

	if *filter == "*" {
		*filter = ".*"
	}

	re := regexp.MustCompile(*filter)

	if *op == "stat" {
		if fi.IsDir() {
			err := cachectl.WalkPrintPagesStat(*fpath, re)
			if err != nil {
				log.Fatalf("failed to walk in %s.", fi.Name())
			}
		} else {
			if !fi.Mode().IsRegular() {
				log.Fatalf("%s is not regular file", fi.Name())
			}

			cachectl.PrintPagesStat(*fpath, fi.Size())
		}
	} else {
		if fi.IsDir() {
			err := cachectl.WalkPurgePages(*fpath, re, *rate, *verbose)
			if err != nil {
				log.Fatalf("failed to walk in %s.", fi.Name())
			}
		} else {
			if !fi.Mode().IsRegular() {
				log.Fatal("%s is not regular file", fi.Name())
			}

			err := cachectl.RunPurgePages(*fpath, fi.Size(), *rate, *verbose)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
