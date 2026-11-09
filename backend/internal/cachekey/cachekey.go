package cachekey

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"time"
)

var salt = newSalt()

func newSalt() []byte {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		binary.LittleEndian.PutUint64(b, uint64(time.Now().UnixNano()))
	}
	return b
}

func Hash(parts ...string) string {
	mac := hmac.New(sha256.New, salt)
	for i, part := range parts {
		if i > 0 {
			mac.Write([]byte{0})
		}
		mac.Write([]byte(part))
	}
	return hex.EncodeToString(mac.Sum(nil))
}
