package chassis

import (
	"math/rand"
)

// Characters to use for random ID generation.
const idChars = "abcdefghjklmnopqrstvwxyzABCDEFGHJKLMNOPQRSTVWXYZ023456789"

const (
	idCharIdxBits = 6                    // 6 bits to represent a character index
	idCharIdxMask = 1<<idCharIdxBits - 1 // All 1-bits, as many as idCharIdxBits
	idCharIdxMax  = 63 / idCharIdxBits   // # of char indices fitting in 63 bits
)

// NewID creates a new random ID.
func NewID(len int) string {
	return RandString(idChars, idCharIdxBits, idCharIdxMask, idCharIdxMax, len)
}

// RandString generates a random string quickly by masking.
// (https://stackoverflow.com/questions/22892120/)
func RandString(chars string, idxBits uint, idxMask int64, idxMax int, n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for idCharIdxMax chars!
	for i, cache, remain := n-1, rand.Int63(), idxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), idxMax
		}
		if idx := int(cache & idxMask); idx < len(chars) {
			b[i] = chars[idx]
			i--
		}
		cache >>= idxBits
		remain--
	}

	return string(b)
}
