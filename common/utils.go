package common

import (
	"io"
	"strings"
)

func readPart(r *strings.Reader, amount int) ([]byte, error) {
	b := make([]byte, amount)
	_, err := r.Read(b)
	return b, err
}

func checkFlag(r *strings.Reader, flag string) (bool, error) {
	byteFlag, err := readPart(r, 1)
	return string(byteFlag) == flag, err
}

func isByteArraysEqual(a1 []byte, a2 []byte) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := 0; i < len(a1); i++ {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func readUntil(r *strings.Reader, marker []byte) []byte {
	var readedData []byte
	for {
		b, err := readPart(r, 1)
		if err == io.EOF {
			break
		}
		readedData = append(readedData, b...)
		if len(readedData) < len(marker) {
			continue
		}
		readedPart := readedData[len(readedData)-len(marker):]
		if isByteArraysEqual(marker, readedPart) {
			break
		}
	}
	return readedData
}
