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
	op := flag.String("op", "stat", "operation(stat, del)")
	fpath := flag.String("f", "", "target file path")
	rate := flag.Float64("r", 1.0, "rate of page cache purged(0.0 <= r<= 1.0)")
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

	if *op == "stat" {
		if fi.IsDir() {
			err := filepath.Walk(*fpath,
				func(path string, info os.FileInfo, err error) error {
					if !info.Mode().IsRegular() {
						return nil
					}
					cachectl.PrintPagesStat(path, info.Size())
					return nil
				})
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
			err := filepath.Walk(*fpath,
				func(path string, info os.FileInfo, err error) error {
					if !info.Mode().IsRegular() {
						return nil
					}
					cachectl.RunPurgePages(path, info.Size(), *rate)
					return nil
				})

			if err != nil {
				fmt.Printf("failed to walk in %s.", fi.Name())
				os.Exit(1)
			}
		} else {
			if !fi.Mode().IsRegular() {
				fmt.Printf("%s is not regular file\n", fi.Name())
				os.Exit(1)
			}

			cachectl.RunPurgePages(*fpath, fi.Size(), *rate)
		}
	}

}
