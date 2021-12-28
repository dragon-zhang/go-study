package demos

import (
	"fmt"
	"golang.org/x/exp/mmap"
	"os"
	"syscall"
)

// see https://geektutu.com/post/quick-go-mmap.html

const maxMapSize = 0x8000000000
const maxMmapStep = 1 << 30 // 1GB

func MmapRead(s string) {
	at, _ := mmap.Open("mmap.txt")
	if at == nil {
		return
	}
	buff := make([]byte, len("Hello world"))
	_, err := at.ReadAt(buff, 0)
	if err != nil {
		fmt.Println(err)
	}
	_ = at.Close()
	fmt.Println(s + " MmapRead->" + string(buff))
}

func mmapSize(size int) (int, error) {
	// Double the size from 32KB until 1GB.
	for i := uint(15); i <= 30; i++ {
		if size <= 1<<i {
			return 1 << i, nil
		}
	}

	// Verify the requested size is not above the maximum allowed.
	if size > maxMapSize {
		return 0, fmt.Errorf("mmap too large")
	}

	// If larger than 1GB then grow by 1GB at a time.
	sz := int64(size)
	if remainder := sz % int64(maxMmapStep); remainder > 0 {
		sz += int64(maxMmapStep) - remainder
	}

	// Ensure that the mmap size is a multiple of the page size.
	// This should always be true since we're incrementing in MBs.
	pageSize := int64(os.Getpagesize())
	if (sz % pageSize) != 0 {
		sz = ((sz / pageSize) + 1) * pageSize
	}

	// If we've exceeded the max size then only grow up to the max size.
	if sz > maxMapSize {
		sz = maxMapSize
	}

	return int(sz), nil
}

func MmapWrite() {
	file, err := os.OpenFile("mmap.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	stat, err := os.Stat("mmap.txt")
	if err != nil {
		panic(err)
	}

	size, err := mmapSize(int(stat.Size()))
	if err != nil {
		panic(err)
	}

	err = syscall.Ftruncate(int(file.Fd()), int64(size))
	if err != nil {
		panic(err)
	}

	buffer, err := syscall.Mmap(int(file.Fd()), 0, size, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	for index, bb := range []byte("Hello world") {
		buffer[index] = bb
	}

	err = syscall.Munmap(buffer)
	if err != nil {
		panic(err)
	}
	fmt.Println("MmapWrite successfully !")
}
