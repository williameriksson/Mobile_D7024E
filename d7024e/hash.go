package d7024e

import (
	"crypto/sha1"
	"io"
	"log"
	"os"
	"encoding/hex"
)

func HashData(data []byte) string {
	h := sha1.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// Hashing a file, probably won't be used
func Hash(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}

// Hashing a string, probably won't be used
func HashStr(text string) []byte {
	h := sha1.New()
	io.WriteString(h, text)
	return h.Sum(nil)
}


