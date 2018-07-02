package servo

import (
	"log"
	"testing"
)

func TestCounterAdd(t *testing.T) {
	var (
		c      *Counter = NewCounter()
		inputs []string = []string{
			"test", "test", "test", "test2", "test3", "test4",
		}
		fn CFn = func(key string, c map[string]int) {
			if c[key] > 2 {
				log.Printf("inside CFn for key(%s)\n", key)
				delete(c, key)
			}
		}
	)
	c.SetCallback(fn)
	for _, v := range inputs {
		c.Add(v)
	}
	if c.Exists("test") {
		t.Fatal("assertion failed, expected equal.")
	}
}

func TestCounterGet(t *testing.T) {
	var (
		c      *Counter = NewCounter()
		inputs []string = []string{
			"test", "test", "test", "test2", "test3", "test4",
		}
		output int
		ok     bool
	)
	for _, v := range inputs {
		c.Add(v)
	}
	output, ok = c.Get("test2")
	if output == -1 || !ok {
		t.Fatal("inconsistent state, unable to get existing value.")
	}
}
