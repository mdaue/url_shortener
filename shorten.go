package main

import (
	"fmt"
	"hash/crc32"
)

func shorten(url string) (string, error) {
	url_bytes := []byte(url)
	hash := crc32.ChecksumIEEE(url_bytes)
	return fmt.Sprintf("%08x", hash), nil
}
