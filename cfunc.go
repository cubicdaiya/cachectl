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
