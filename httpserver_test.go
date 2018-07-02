package servo

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

type assertOpts struct {
	opts   HTTPServerOpts
	assert bool
}

func TestOpts(t *testing.T) {
	var (
		testcases []assertOpts
		ok        bool
		err       error
	)
	testcases = []assertOpts{
		assertOpts{
			opts: HTTPServerOpts{
				Address:  "",
				Port:     0,
				MaxConns: 0,
			},
			assert: true,
		},
		assertOpts{
			opts: HTTPServerOpts{
				Address:  "loaclhost",
				Port:     8080,
				MaxConns: 1024,
			},
			assert: false,
		},
		assertOpts{
			opts: HTTPServerOpts{
				Address:  "192.168.1.7",
				Port:     9999,
				MaxConns: 4096,
			},
			assert: false,
		},
		assertOpts{
			opts: HTTPServerOpts{
				Address:  "0.0.0.0",
				Port:     80,
				MaxConns: 2048,
			},
			assert: true,
		},
		assertOpts{
			opts: HTTPServerOpts{
				Address:  "test-addr.com",
				Port:     9090,
				MaxConns: 0,
			},
			assert: false,
		},
		assertOpts{
			opts: HTTPServerOpts{
				Address:  "test-addr.com",
				Port:     75535,
				MaxConns: 0,
			},
			assert: false,
		},
	}
	for i, tcase := range testcases {
		ok, err = checkAndSetOpts(tcase.opts)
		if err != nil && ok != tcase.assert {
			t.Fatalf("- [%d] inconsistent state, expected equal. ( \nerr=%v\n case=%v\n ok=%t\n assert=%t\n) ", i, err, tcase, ok, tcase.assert)
		}
	}
}

func TestRoutes(t *testing.T) {
	var (
		r       *Routes = NewRoutes()
		handler HandlerFn
		ok      bool
	)
	handler = func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("routes: received the request.")
	}
	// add a route to Rouets struct
	r.Register("/", handler)
	ok = r.Exists("/")
	if !ok {
		t.Fatal("assertion failed, expected true for existing path.")
	}
	ok = r.Remove("/")
	if !ok {
		t.Fatal("assertion failed, expected true for existing path.")
	}
}

// - MARK: StatServer section.

func nHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("received a new reuqest")
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("inside the main handler.")
}

func TestServer(t *testing.T) {
	const (
		cDEFAULTTIMEOUT time.Duration = time.Duration(time.Second * 2)
	)
	var (
		s     *HTTPServer     = NewHTTPServer()
		wg    *sync.WaitGroup = &sync.WaitGroup{}
		sopts HTTPServerOpts
		err   error
		ok    bool
	)
	sopts = HTTPServerOpts{
		Address:  "localhost",
		Port:     8080,
		MaxConns: 100,
	}
	// START - refactor
	s.SetLogger(nil)
	// END - refactor
	err = s.Setup(sopts)
	if err != nil {
		t.Fatal("assertion failed, expected nil value.", err)
	}
	if s.Logger == nil {
		t.Fatal("inconsistent state, logger is null.")
	}
	_ = s.Logger.Write("output")
	s.Register("/", mainHandler)
	s.Register("/main", nHandler)
	// it panics when suplied mux is empty
	// ok = s.routes.SubmitToMux(s.mux)
	ok = s.SubmitRoutes()
	if !ok {
		t.Fatal("inconsistent state, expected true value.")
	}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		err = s.Start()
		fmt.Println("error from goroutine: ->", err)
		wg.Done()
	}(wg)
	time.Sleep(cDEFAULTTIMEOUT)
	s.Shutdown()
	wg.Wait()
	fmt.Println("wait gorup is done")
}
