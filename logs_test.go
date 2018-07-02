package servo

import (
	"log"
	"os"
	"testing"
)

func TestNewLog(t *testing.T) {
	var (
		l   *LogFile = NewLog()
		err error
	)
	l.opts = LogOpts{
		output: "/tmp/f.txt",
	}
	err = l.Setup()
	if err != nil {
		t.Fatal("inconsistent state, expected err == nil.", err)
	}
	// log.Println("output")
}

func TestRAWLog(t *testing.T) {
	var (
		file *os.File
		err  error
		l    *log.Logger
	)
	file, err = os.OpenFile("/tmp/test.txt", os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		t.Fatal("unable to open/create the file.")
	}
	l = log.New(file, "servo", log.LstdFlags) // standard OS flags
	l.SetOutput(file)
	l.Println("testoutput")
}
