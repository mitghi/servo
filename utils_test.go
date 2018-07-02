package servo

import (
	"fmt"
	"testing"
)

func TestRandString(t *testing.T) {
	var (
		s   string
		err error
	)
	s, err = RandString16()
	if err != nil {
		t.Fatal("assertion failed, expected nil.")
	}
	fmt.Println(s)
}
