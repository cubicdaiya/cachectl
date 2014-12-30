package cachectl

/*
#include <stdlib.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>

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

*/
import "C"
import (
	"errors"
	"fmt"
)

func purgePages(fpath string, fsize int64, rate float64) error {
	if rate < 0.0 || rate > 1.0 {
		return errors.New(fmt.Sprintf("%f: rate should be less than 1.0\n", rate))
	}

	result := C.fadvise(C.CString(fpath), C.float(rate))
	if result == -1 {
		return errors.New(fmt.Sprintf("failed to purge page cache for %s", fpath))
	}

	return nil
}

func RunPurgePages(path string, fsize int64, rate float64) {
	fmt.Printf("Before purging %s 's page cache\n\n", path)
	PrintPagesStat(path, fsize)

	purgePages(path, fsize, rate)

	fmt.Printf("\nAfter purging %s 's page cache\n\n", path)
	PrintPagesStat(path, fsize)
	fmt.Println()
}
