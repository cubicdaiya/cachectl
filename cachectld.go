package main

import (
	"./cachectl"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func scheduledPurgePages(target cachectl.SectionTarget) {

	re := regexp.MustCompile(target.Filter)

	for {
		timer := time.NewTimer(time.Second * time.Duration(target.PurgeInterval))
		<-timer.C

		fi, err := os.Stat(target.Path)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		if fi.IsDir() {
			err := cachectl.WalkPurgePages(target.Path, re, target.Rate)
			if err != nil {
				fmt.Printf("failed to walk in %s.", fi.Name())
			}
		} else {
			if !fi.Mode().IsRegular() {
				fmt.Printf("%s is not regular file\n", fi.Name())
				continue
			}

			cachectl.RunPurgePages(target.Path, fi.Size(), target.Rate)
		}
	}
}

func waitSignal() int {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var exitcode int

	s := <-sigchan

	switch s {
	case syscall.SIGHUP:
		fallthrough
	case syscall.SIGINT:
		fallthrough
	case syscall.SIGTERM:
		fallthrough
	case syscall.SIGQUIT:
		exitcode = 0
	default:
		exitcode = 1
	}

	return exitcode
}

func main() {

	// Parse flags
	version := flag.Bool("v", false, "show version")
	confPath := flag.String("c", "", "configuration file for cachectld")
	flag.Parse()

	if *version {
		cachectl.PrintVersion(cachectl.Cachectld)
		os.Exit(0)
	}

	var confCachectld cachectl.ConfToml
	err := cachectl.LoadConf(*confPath, &confCachectld)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cachectl.ValidateConf(&confCachectld)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, target := range confCachectld.Targets {
		go scheduledPurgePages(target)
	}

	code := waitSignal()

	os.Exit(code)
}
