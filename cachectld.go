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

func purgePages(target cachectl.SectionTarget, re *regexp.Regexp) error {
	fi, err := os.Stat(target.Path)
	if err != nil {
		return err
	}

	verbose := false

	if fi.IsDir() {
		err := cachectl.WalkPurgePages(target.Path, re, target.Rate, verbose)
		if err != nil {
			return fmt.Errorf("failed to walk in %s.", fi.Name())
		}
	} else {
		if !fi.Mode().IsRegular() {
			return fmt.Errorf("%s is not regular file", fi.Name())
		}

		err := cachectl.RunPurgePages(target.Path, fi.Size(), target.Rate, verbose)
		if err != nil {
			return fmt.Errorf("%s: %s", fi.Name(), err.Error())
		}
	}

	return nil
}

func scheduledPurgePages(target cachectl.SectionTarget) {

	if target.PurgeInterval == -1 {
		log.Printf("cachectld runs for the target(path:%s, filter:%s) when only received USR1\n",
			target.Path, target.Filter)
		return
	}

	re := regexp.MustCompile(target.Filter)

	for {
		timer := time.NewTimer(time.Second * time.Duration(target.PurgeInterval))
		<-timer.C

		err := purgePages(target, re)
		if err != nil {
			log.Println(err)
		}
	}
}

func waitSignal() int {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan,
		syscall.SIGUSR1,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	var exitcode int

	s := <-sigchan

	switch s {
	case syscall.SIGUSR1:
		// not exit
		exitcode = -1
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
		log.Fatal(err)
	}

	err = cachectl.ValidateConf(&confCachectld)
	if err != nil {
		log.Fatal(err)
	}

	for _, target := range confCachectld.Targets {
		go scheduledPurgePages(target)
	}

waitSignalLoop:
	code := waitSignal()

	// When received SIGUSR1,
	// cachectld runs purgePages().
	if code == -1 {
		for _, target := range confCachectld.Targets {
			re := regexp.MustCompile(target.Filter)
			err := purgePages(target, re)
			if err != nil {
				log.Println(err)
			}
		}
		goto waitSignalLoop
	}

	os.Exit(code)
}
