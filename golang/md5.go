package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
)

func main() {

	data := []byte("Hello World")
	fmt.Printf("%x\n", md5.Sum(data))
	fmt.Printf("%x\n", crc32.ChecksumIEEE(data))
}

func Crc32IEEE(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func Md5(data []byte) [16]byte {
	return md5.Sum(data)
}
