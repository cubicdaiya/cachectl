package cachectl

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/sys/unix"
)

func purgePages(fpath string, fsize int64, rate float64) error {
	if rate < 0.0 || rate > 1.0 {
		return fmt.Errorf("%.1f: rate should be over 0.0 and less than 1.0\n", rate)
	}

	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = unix.Fadvise(int(f.Fd()), 0, int64(float64(fsize)*rate), unix.FADV_DONTNEED)
	if err != nil {
		return fmt.Errorf("failed to purge page cache for %s", fpath)
	}

	return nil
}

func RunPurgePages(path string, fsize int64, rate float64, verbose bool) error {
	if verbose {
		fmt.Printf("Before purging %s 's page cache\n\n", path)
		PrintPagesStat(path, fsize)
	} else {
		log.Printf("purging %s 's page cache\n", path)
	}

	err := purgePages(path, fsize, rate)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("\nAfter purging %s 's page cache\n\n", path)
		PrintPagesStat(path, fsize)
		fmt.Println()
	}

	return nil
}
