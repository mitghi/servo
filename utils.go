package servo

import (
	"encoding/base64"
	"log"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	// cChars is default ASCII translation table
	cChars []byte = []byte{0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a}
)

// Rbtosftable is a function that generates a random string by
// translating random bytes based on `table`. It returns an
// error in case of unsuccessfull operation.
func Rbtosftable(size uint, table []byte) (string, error) {
	var (
		tl  int    = len(table)       // table length
		mch int    = 255 - (256 % tl) // max character code allowed
		l   int    = int(size)        // size as int
		pos int    = 0                // cursor position
		err error  = EINVAL           // default error message
		rb  []byte                    // random bytes
	)
	if size == 0 || (tl < 2 || tl > 256) {
		return _EMPTY_, EINVAL
	}
	rb = make([]byte, l)
	for _, err = rand.Read(rb); err == nil; {
		for i := 0; i < l; i++ {
			var c int = int(rb[i])
			if c <= mch {
				rb[pos] = table[c%tl]
				pos++
				if pos == l {
					return string(rb), nil
				}
			}
		}
	}

	return _EMPTY_, err
}

// RandBytes is a function that returns `size` random bytes.
// It returns an error in case of failure.
func RandBytes(size uint) ([]byte, error) {
	var (
		buff []byte = make([]byte, size) // rand byte buffer
		err  error
	)
	_, err = rand.Read(buff)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

// RandString returns a random `size` string. It returns
// an error in case of failure. Note; the underlaying
// hashing mechanism used by this function is `base64`.
func RandString(size uint) (r string, err error) {
	var (
		b []byte
	)
	b, err = RandBytes(size)
	if err != nil {
		return _EMPTY_, err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RandString16 is a `Rtosftable` wrapper function that
// returns random string of size 16.
func RandString16() (string, error) {
	return Rbtosftable(16, cChars)
}

// RandString16 is a `Rtosftable` wrapper function that
// returns random string of size 20 ( UUID length ).
func RandString20() (string, error) {
	return Rbtosftable(20, cChars)
}

func FOpenOrCreate(path string) (f *os.File) {
	var (
		err error
	)
	f, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0660)
	if err == nil {
		return f
	}
	f, err = os.Create(path)
	if err != nil {
		// TODO:
		// . refactor this
		log.Println("(FOpenOrCreate) unable to create the file.")
		return nil
	}
	return f
}
