package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
)

const CACHECTL_VERSION = "0.0.1"

func printCachectlVersion() {
	fmt.Printf(`cachectl %s
Compiler: %s %s
Copyright (C) 2014 Tatsuhiko Kubo <cubicdaiya@gmail.com>
`,
		CACHECTL_VERSION,
		runtime.Compiler,
		runtime.Version())
}

func printCacheStat(fpath string, fsize int64) {
	pagesize := os.Getpagesize()
	pagesizeKB := pagesize / 1024
	if fsize == 0 {
		fmt.Printf("%s 's pages in cache: %d/%d (%.1f%%) [filesize=%.1fK, pagesize=%dK]\n", fpath, 0, 0, 0.0, 0.0, pagesizeKB)
	}

	pages := (fsize + int64(pagesize) - 1) / int64(pagesize)
	pagesActive := C.activePages(C.CString(fpath))
	activeRate := 100.0 * (float64(pagesActive) / float64(pages))
	filesizeKB := float64(fsize) / 1024
	fmt.Printf("%s 's pages in cache: %d/%d (%.1f%%)  [filesize=%.1fK, pagesize=%dK]\n",
		fpath, pagesActive, pages, activeRate, filesizeKB, pagesizeKB)
}

func deleteCache(fpath string, fsize int64, rate float64) error {
	if rate < 0.0 || rate > 1.0 {
		return errors.New(fmt.Sprintf("%f: rate should be less than 1.0\n", rate))
	}

	result := C.fadvise(C.CString(fpath), C.float(rate))
	if result == -1 {
		return errors.New(fmt.Sprintf("failed to delete page cache for %s", fpath))
	}

	return nil
}

func main() {

	// Parse flags
	version := flag.Bool("v", false, "show version")
	op := flag.String("op", "stat", "operation(stat, del)")
	fpath := flag.String("f", "", "target file path")
	rate := flag.Float64("r", 1.0, "rate of page cache deleted(0.0 <= r<= 1.0)")
	flag.Parse()

	if *version {
		printCachectlVersion()
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

	if !fi.Mode().IsRegular() {
		fmt.Printf("%s is not regular file\n", fi.Name())
		os.Exit(1)
	}

	if *op == "stat" {
		printCacheStat(*fpath, fi.Size())
	} else {
		fmt.Printf("Before deleting %s 's page cache\n\n", *fpath)
		printCacheStat(*fpath, fi.Size())

		deleteCache(*fpath, fi.Size(), *rate)

		fmt.Printf("\nAfter deleting %s 's page cache\n\n", *fpath)
		printCacheStat(*fpath, fi.Size())
	}

}
