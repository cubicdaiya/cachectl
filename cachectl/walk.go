package cachectl

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func WalkPrintPagesStat(path string, re *regexp.Regexp) error {
	return filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err.Error())
				return nil
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			if re.MatchString(path) {
				PrintPagesStat(path, info.Size())
			}
			return nil
		})
}

func WalkPurgePages(path string, re *regexp.Regexp, rate float64, verbose bool) error {
	return filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Println(err.Error())
				return nil
			}
			if !info.Mode().IsRegular() {
				return nil
			}

			if re.MatchString(path) {
				err := RunPurgePages(path, info.Size(), rate, verbose)
				if err != nil {
					log.Println(err.Error())
				}
			}
			return nil
		})
}
