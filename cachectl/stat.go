package cachectl

import (
	"log"
	"os"
)

func PrintPagesStat(fpath string, fsize int64) {
	pagesize := os.Getpagesize()
	pagesizeKB := pagesize / 1024
	if fsize == 0 {
		log.Printf("%s 's pages in cache: %d/%d (%.1f%%) [filesize=%.1fK, pagesize=%dK]\n", fpath, 0, 0, 0.0, 0.0, pagesizeKB)
		return
	}

	pages := (fsize + int64(pagesize) - 1) / int64(pagesize)

	pagesActive := activePages(fpath)
	activeRate := float64(0)
	if pagesActive == -1 {
		pagesActive = 0
		pages = 0
	} else {
		activeRate = 100.0 * (float64(pagesActive) / float64(pages))
	}
	filesizeKB := float64(fsize) / 1024
	log.Printf("%s 's pages in cache: %d/%d (%.1f%%)  [filesize=%.1fK, pagesize=%dK]\n",
		fpath, pagesActive, pages, activeRate, filesizeKB, pagesizeKB)
}
