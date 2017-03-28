package cachectl

import (
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func activePages(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	fsize := fi.Size()

	if fsize == 0 {
		return 0, nil
	}

	mmap, err := unix.Mmap(int(f.Fd()), 0, int(fsize), unix.PROT_NONE, unix.MAP_SHARED)
	if err != nil {
		return 0, err
	}

	pagesize := int64(os.Getpagesize())
	pages := (fsize + pagesize - 1) / pagesize
	pageinfo := make([]byte, pages)

	mmapPtr := uintptr(unsafe.Pointer(&mmap[0]))
	sizePtr := uintptr(fsize)
	pageinfoPtr := uintptr(unsafe.Pointer(&pageinfo[0]))

	ret, _, err := unix.Syscall(unix.SYS_MINCORE, mmapPtr, sizePtr, pageinfoPtr)
	if ret != 0 {
		return 0, err
	}
	defer unix.Munmap(mmap)

	result := 0

	for _, p := range pageinfo {
		if p%2 == 1 {
			result++
		}
	}

	return result, nil
}
