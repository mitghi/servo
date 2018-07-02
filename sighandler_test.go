package servo

import (
	"log"
	"os"
	"syscall"
	"testing"
)

type tstsigfn func(os.Signal)

func tstHandleSIGHUP(sig os.Signal) {
	log.Println("(tstHandleSIGHUP) handler invoked.")
}

func TestSignalRegisteration(t *testing.T) {
	var (
		s     *HTTPServer = NewHTTPServer()
		sopts HTTPServerOpts
		err   error
		ok    bool
	)
	sopts = HTTPServerOpts{
		Address:  "localhost",
		Port:     8080,
		MaxConns: 1000,
	}
	// perform inital state initialization
	err = s.Setup(sopts)
	if err != nil {
		t.Fatal("inconsistent state, unable to setup httpserver instance.")
	}
	// register signal handlers ( NOTE:
	// signal handlers can be used to
	// customize program's behavior
	// when a certain signal is catched.
	s.RegisterSignalHandler(syscall.SIGHUP, tstHandleSIGHUP)
	/* d e b u g */
	// if !{
	// 	t.Fatal("assertion failed, expected true.")
	// }
	// if {
	// 	t.Fatal("assertion failed, expected false.")
	// }
	/* d e b u g */
	// commit handlers to dispatcher
	s.CommitSignalHandlers()
	// test (*HTTPServer) dispatchSig
	// test `sigContainer` :sigbox
	// NOTE:
	// . sigs become usable after `CommitSignalHandlers`.
	ok = s.sigbox.sigs.Exists(syscall.SIGHUP)
	if !ok {
		t.Fatal("assertion failed, expected true. Unable to retrieve existing handler.")
	}
	err = s.dispatchSig(syscall.SIGHUP)
	if err != nil {
		t.Fatalf("inconsistent state, expected err==nil. got : %+v\n", err)
	}
}
