package shorten

import (
	"hash/crc64"
	"strings"
)

const (
	alphabet       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	alphabetLength = 63
)

// Хэширование url
func HashString(s string) (string, error) {

	// Hashing given string with CRC-64
	table := crc64.MakeTable(crc64.ISO)
	hash := crc64.Checksum([]byte(s), table)
	array := make([]uint8, 0, 10)
	for mod := uint64(0); hash != 0 && len(array) < 10; {
		mod = hash % alphabetLength
		hash /= alphabetLength
		array = append(array, alphabet[mod])
	}
	return strings.Repeat(string(alphabet[0]), 10-len(array)) + string(array), nil
}
