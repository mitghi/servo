package servo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/mitghi/util/atoms"
)

const (
	cSSFmt     string = "statserver: [%-8s] %s"
	cPUBADDR   string = "0.0.0.0"
	cLHOSTADDR string = "localhost"
	cHOMEADDR  string = "127.0.0.1"
)

const (
	defVALIDURILEN int    = 4
	defLogOutput   string = "/tmp/servo.log"
)

var (
	ESTATSRVINCOMP     error = errors.New("statserver: incomplete operation/request")
	ESTATSRVINVAL      error = errors.New("statserver: invalid operation")
	ESTATSRVFATAL      error = errors.New("statserver: fatal condition encountered")
	ESTATSRVNOTRUNNING error = errors.New("Server: server instance is not running.")
	ESTATSRVRUNNING    error = errors.New("Server: server already in running mode.")
	ESRVSIGFATAL       error = errors.New("HTTPServer: failed to commit signal handlers.")
	ESRVINVKFATAL      error = errors.New("HTTPServer: Failed to retrieve/invoke handler.")
	ESRVLKUPFATAL      error = errors.New("HTTPServer: Unable to find appropirate sig handler.")
)

// - MARK: Utilities section.
func debugPrint(args ...interface{}) {
	log.Printf(cSSFmt, args...)
}

func dbgpFatal(input string) {
	fmt.Printf("statserver: [%s] %s.\n", "fatal", input)
}

// - MARK: Alloc/Init section.

func NewHTTPServer() (s *HTTPServer) {
	s = &HTTPServer{
		exitCh:   make(chan os.Signal, 1),
		running:  atoms.NewBoolean(),
		doneInit: atoms.NewBoolean(),
		Varz:     &ServerVarz{},
		mu:       &sync.RWMutex{},
		routes:   NewRoutes(),
		server:   nil,
		sigbox:   NewSigContainer(),
	}
	return s
}

func NewRoutes() *Routes {
	return &Routes{
		mu:     &sync.RWMutex{},
		routes: make(map[string]HandlerFn),
	}
}

// - MARK: HTTPServer section.

func (s *HTTPServer) hasValidState() (ok bool, err error) {
	/* d e b u g */
	// if s.parent == nil {
	// 	err = ESTATSRVINCOMP
	// 	goto ERROR
	// }
	/* d e b u g */
	if s.mux != nil {
		err = ESTATSRVFATAL
		goto ERROR
	}
	if s.mux != nil {
		err = ESTATSRVFATAL
		goto ERROR
	}
	return true, nil
ERROR:
	return false, err
}

func (s *HTTPServer) Setup(opts HTTPServerOpts) (err error) {
	var (
		ok bool
	)
	if s.doneInit.Is(true) {
		// reinvoking `Setup(...)` in fully
		// initialized state
		return ESTATSRVFATAL
	}
	// TODO:
	// . refactor locking mechs.
	// s.mu.Lock()
	// defer s.mu.Unlock()
	if ok, err = s.hasValidState(); !ok || err != nil {
		goto ERROR
	}
	if ok, err = s.checkAndSetOpts(opts, true); err != nil {
		goto ERROR
	}
	if s.server != nil {
		log.Println("server: invalid state, expected http server to not be initialized.")
	}
	/* d e b u g */
	// fmt.Println("address is:", s.addr)
	/* d e b u g */
	s.server = &http.Server{
		Addr:    s.addr,
		Handler: s,
	}
	s.mux = http.NewServeMux()
	s.DoneCh = make(chan struct{})
	s.hasSignals = atoms.NewBoolean()
	if len(opts.Handlers) != 0 {
		/* d e b u g */
		// output number of registered handlers
		log.Println("Handlers:", len(opts.Handlers))
		/* d e b u g */
		for _, route := range opts.Handlers {
			s.Register(route.Path, route.Handler)
		}
	}
	// TODO:
	// . refactor this
	ok = s.SubmitRoutes()
	if !ok {
		err = ESRVFATALOPS
		log.Println("Fatal cond while submiting routes.")
		goto ERROR
	}
	// TODO:
	// . review this
	// s.Logger = NewLogFromOpts(LogOpts{
	// 	output: cDefaultLogPath,
	// })
	s.Logger = GetDefault()
	s.doneInit.Set(true)
	return nil
ERROR:
	return err
	/* d e b u g */
	// CLEANUP:
	// 	if s.server != nil {
	// 		s.server = nil
	// 	}
	// 	if s.mux != nil {
	// 		s.mux = nil
	// 	}
	// 	if s.DoneCh != nil {
	// 		s.DoneCh = nil
	// 	}
	/* d e b u g */
}

func (s *HTTPServer) checkAndSetOpts(opts HTTPServerOpts, shouldSet bool) (bool, error) {
	var (
		// isValid assumes `opts` is valid ( by default )
		// unless proven otherwise.
		isValid bool = true
		host    string
		port    string
		ok      bool
		err     error
	)
	switch opts.Address {
	case "", cLHOSTADDR, cHOMEADDR, cPUBADDR:
	default:
		host, _, err = net.SplitHostPort(opts.Address)
		if err != nil || !ok {
			dbgpFatal("invalid address")
			goto ERROR
		}
		if host != opts.Address {
			dbgpFatal("invalid host")
			err = ESTATSRVINVAL
			goto ERROR
		}
	}
	if opts.Port < 0 || opts.Port > 65535 {
		dbgpFatal("invalid port")
		isValid = false
		err = ESTATSRVINVAL
		goto ERROR
	}
	if opts.MaxConns < 0 {
		dbgpFatal("invalid maxconns")
		isValid = false
		err = ESTATSRVINVAL
		goto ERROR
	}
	if shouldSet {
		// set listener address
		port = strconv.Itoa(opts.Port)
		s.addr = strings.Join([]string{opts.Address, port}, ":")
	}
	return isValid, nil
ERROR:
	return false, err
}

func (s *HTTPServer) Start() (err error) {
	if s.running.Is(true) {
		log.Println("-> fatal 1")
		err = ESRVFATALOPS
		goto ERROR
	} else if s.doneInit.Is(false) {
		log.Println("-> fatal 2")
		err = ESRVFATALOPS
		goto ERROR
	}
	s.running.Set(true)
	err = s.server.ListenAndServe()
	if err != nil {
		log.Println("- (HTTPServer) unable to start http server.", err)
	}
	/* d e b u g */
	// log.Println("error from starting server: ", err)
	/* d e b u g */
ERROR:
	return err
}

func (s *HTTPServer) Shutdown() (err error) {
	var (
		ctx context.Context
	)
	if s.running.Is(false) {
		return ESTATSRVNOTRUNNING
	}
	ctx = context.Background()
	err = s.server.Shutdown(ctx)
	/* d e b u g */
	// log.Println("error from StatsServer:->", err)
	/* d e b u g */
	if err != nil {
		s.DoneCh <- struct{}{}
		s.running.Set(false)
		return nil
	}
	return err
}

func (s *HTTPServer) SubmitRoutes() (ok bool) {
	s.mu.RLock()
	if s.mux == nil {
		ok = false
		s.mu.RUnlock()
		return false
	}
	ok = s.routes.SubmitToMux(s.mux)
	s.mu.RUnlock()
	return ok
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		handler http.Handler
		// pattern string
	)
	/* d e b u g */
	// log.Println("handler & pattern", handler, pattern)
	// log.Printf("(HTTPServer/ServeHTTP): mux handler with pattern(%s).\n", pattern)
	/* d e b u g */
	handler, _ = s.mux.Handler(r)
	handler.ServeHTTP(w, r)
	s.Varz.Inc(Connection)
}

func (s *HTTPServer) Register(path string, handler HandlerFn) {
	s.mu.Lock()
	s.routes.Register(path, handler)
	s.mu.Unlock()
}

// Wait blocks on the exit channel until server
// either shuts down or terminates.
func (s *HTTPServer) Wait() {
	// TODO:
	// . null check for exit channel
	if s.running.Is(true) {
		<-s.DoneCh
	}
	return
}

func (s *HTTPServer) MarkDone() {
	// TODO:
	// . null check for exit channel
	if s.running.Is(true) {
		s.DoneCh <- struct{}{}
	}
	return
}

// - MARK: Utility section.

func checkAndSetOpts(opts HTTPServerOpts) (bool, error) {
	var (
		// isValid assumes `opts` is valid ( by default )
		// unless proven otherwise.
		isValid bool = true
		host    string
		ok      bool
		err     error
	)
	switch opts.Address {
	case "", cLHOSTADDR, cHOMEADDR, cPUBADDR:
	default:
		host, _, err = net.SplitHostPort(opts.Address)
		if err != nil || !ok {
			dbgpFatal("invalid address")
			goto ERROR
		}
		if host != opts.Address {
			dbgpFatal("invalid host")
			err = ESTATSRVINVAL
			goto ERROR
		}
	}
	if opts.Port < 0 || opts.Port > 65535 {
		dbgpFatal("invalid port")
		isValid = false
		err = ESTATSRVINVAL
		goto ERROR
	}
	if opts.MaxConns < 0 {
		dbgpFatal("invalid maxconns")
		isValid = false
		err = ESTATSRVINVAL
		goto ERROR
	}
	return isValid, nil
ERROR:
	return false, err
}

func (s *HTTPServer) CommitSignalHandlers() (err error) {
	// NOTE
	// . refactor this receiver method into
	//   interface ( protocol ) conformance.
	var (
		ok bool
	)
	// s.mu.Lock()
	// s.mu.Unlock()
	log.Println("s.sigbox->", s.sigbox == nil, s.sigbox)
	ok = s.sigbox.Commit()
	if !ok {
		// TODO:
		return ESRVSIGFATAL
	}
	return nil
}

func (s *HTTPServer) SetLogger(l *LogFile) (ok bool) {
	var (
		logger *LogFile
		err    error
	)
	s.mu.Lock()
	if s.Logger == nil {
		logger = NewLogFromOpts(LogOpts{
			output: cDefaultLogPath,
		})
		err = logger.Setup()
		if err != nil {
			log.Println("- (HTTPServer) unable to setup logger.", err)
			s.mu.Unlock()
			return false
		}
	}
	s.Logger = logger
	s.mu.Unlock()
	return true
}

func (s *HTTPServer) Address() string {
	return s.server.Addr
}
