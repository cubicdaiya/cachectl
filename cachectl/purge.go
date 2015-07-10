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
	"fmt"
	"log"
	"unsafe"
)

func purgePages(fpath string, fsize int64, rate float64) error {
	if rate < 0.0 || rate > 1.0 {
		return fmt.Errorf("%.1f: rate should be less than 1.0\n", rate)
	}

	cs := C.CString(fpath)
	defer C.free(unsafe.Pointer(cs))
	result := C.fadvise(cs, C.float(rate))
	if result == -1 {
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
