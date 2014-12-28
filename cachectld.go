package main

import (
	"./cachectl"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

func ScheduledPurgePages(target cachectl.SectionTarget) {

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
			err := filepath.Walk(target.Path,
				func(path string, info os.FileInfo, err error) error {
					if !info.Mode().IsRegular() {
						return nil
					}

					if re.MatchString(path) {
						cachectl.RunPurgePages(path, info.Size(), target.Rate)
					}
					return nil
				})

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
		panic(err.Error())
	}

	err = cachectl.ValidateConf(&confCachectld)
	if err != nil {
		panic(err.Error())
	}

	for _, target := range confCachectld.Targets {
		go ScheduledPurgePages(target)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	exitchan := make(chan int)
	go func() {
		for {
			s := <-sigchan
			switch s {
			case syscall.SIGHUP:
				fallthrough
			case syscall.SIGINT:
				fallthrough
			case syscall.SIGTERM:
				fallthrough
			case syscall.SIGQUIT:
				exitchan <- 0
			default:
				exitchan <- 1
			}
		}
	}()

	code := <-exitchan
	os.Exit(code)
}
