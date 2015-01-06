package main

import (
	"./cachectl"
	"flag"
	"fmt"
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
	flag.Parse()

	if *version {
		cachectl.PrintVersion(cachectl.Cachectl)
		os.Exit(0)
	}

	if *fpath == "" {
		fmt.Println("target file path is empty.")
		os.Exit(0)
	}

	fi, err := os.Stat(*fpath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *filter == "*" {
		*filter = ".*"
	}

	re := regexp.MustCompile(*filter)
	verbose := true

	if *op == "stat" {
		if fi.IsDir() {
			err := cachectl.WalkPrintPagesStat(*fpath, re)
			if err != nil {
				fmt.Printf("failed to walk in %s.", fi.Name())
				os.Exit(1)
			}
		} else {
			if !fi.Mode().IsRegular() {
				fmt.Printf("%s is not regular file\n", fi.Name())
				os.Exit(1)
			}

			cachectl.PrintPagesStat(*fpath, fi.Size())
		}
	} else {
		if fi.IsDir() {
			err := cachectl.WalkPurgePages(*fpath, re, *rate, verbose)
			if err != nil {
				fmt.Printf("failed to walk in %s.", fi.Name())
				os.Exit(1)
			}
		} else {
			if !fi.Mode().IsRegular() {
				fmt.Printf("%s is not regular file\n", fi.Name())
				os.Exit(1)
			}

			cachectl.RunPurgePages(*fpath, fi.Size(), *rate, verbose)
		}
	}
}
