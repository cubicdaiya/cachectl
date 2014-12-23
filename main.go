package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/mman.h>

int fadvise(const char *path, float r)
{
    int fd;
    struct stat st;
    off_t l;
    fd = open(path, O_RDONLY);
    if(fd == -1) {
        return -1;
    }

    if(fstat(fd, &st) == -1) {
        goto error;
    }

    l = (off_t)(st.st_size * r);

    if(posix_fadvise(fd, 0, l, POSIX_FADV_DONTNEED) != 0) {
        goto error;
    }

    close(fd);
    return 1;
error:
    close(fd);
    return -1;
}

int activePages(const char *path)
{
    int i, j, fd, pages, pagesize;
    struct stat st;
    void *m;
    char *pageinfo;

    fd = open(path, O_RDONLY);
    if(fd == -1) {
        return -1;
    }

    if(fstat(fd, &st) == -1) {
        goto error;
    }

    pagesize = getpagesize();
    pages = (st.st_size + pagesize - 1) / pagesize;
    pageinfo = calloc(sizeof(*pageinfo), pages);
    if(!pageinfo) {
        goto error;
    }

    m = mmap(NULL, st.st_size, PROT_NONE, MAP_SHARED, fd, 0);
    if(m == MAP_FAILED) {
        free(pageinfo);
        goto error;
    }

    if(mincore(m, st.st_size, pageinfo) == -1) {
        free(pageinfo);
        munmap(m, st.st_size);
        goto error;
    }

    i = 0;
    j = 0;
    for (i = 0; i < pages; i++) {
        if(pageinfo[i++] & 1) {
            j++;
        }
    }

    munmap(m, st.st_size);

    return j;
error:
    close(fd);
    return -1;
}
*/
import "C"
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

	result := C.fadvise(fpath, rate)
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
		fmt.Printf("Before deleting %s 's page cache\n", *fpath)
		printCacheStat(*fpath, fi.Size())

		deleteCache(*fpath, fi.Size(), *rate)

		fmt.Printf("After deleting %s 's page cache\n", *fpath)
		printCacheStat(*fpath, fi.Size())
	}

}
