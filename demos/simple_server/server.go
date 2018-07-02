package main

import (
	"fmt"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/mitghi/servo"
)

// TODO:
// . fetch statistics via SIGHUP

// - MARK: Constant section.

const (
	cDefaultOutputPath string = "/tmp"
)

// Flags
var (
	ip      *string
	port    *int
	throtle *int
)

// Server is the container for handling incoming requests.
type Server struct {
	*servo.HTTPServer
	mu      *sync.Mutex
	Counter *servo.Counter
}

// Globals
var (
	s      *Server			// serving entry
	sopts  servo.HTTPServerOpts	// serving configurations
	exitCh chan os.Signal		// termination channel; for graceful shutdown
	routes []servo.Route		// routing table
	err    error
)

// - MARK: Alloc/Init section.

func NewServer() *Server {
	return &Server{
		HTTPServer: servo.NewHTTPServer(),
		mu:         &sync.Mutex{},
		Counter:    servo.NewCounter(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	s.HTTPServer.ServeHTTP(w, r)
}

func createPath(args ...string) string {
	return strings.Join(args, "/")
}

func counterCB(key string, c map[string]int) {
	var (
		path string
		err  error
	)
	log.Printf("(counterCB) is invoked, key: (%s)\n", key)
	log.Printf("(counterCB) key_length: (%d)\n", c[key])
	if c[key] > 2 {
		// NOTE
		// . key must not be purged
		//   in order to remain
		//   in consistent state.
		path = createPath(cDefaultOutputPath, key)
		if _, err = os.Stat(path); err == nil {
			/* d e b u g */
			log.Printf("(counterCB) stat : %+v\n", err)
			/* d e b u g */
			if err = os.Remove(path); err != nil {
				log.Println("(counterCB) cannot remove file.", err)
			}
		}
	}
}

func main() {
	// initialize server instance
	s = NewServer()
	s.Counter.SetCallback(counterCB)
	s.Counter.SetThreshold(2)
	exitCh = make(chan os.Signal, 1)
	// register for exit signals
	signal.Notify(exitCh, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	// setup flags
	ip = flag.String("address", "127.0.0.1", "serving ip")
	port = flag.Int("port", 8080, "serving port")
	throtle = flag.Int("throtle", 100, "throtle after :n maximum connections")
	flag.Parse()
	// setup route definitions
	routes = []servo.Route{
		servo.Route{
			Path:    "/",
			Handler: s.indexHandler,
		},
		servo.Route{
			Path:    "/main",
			Handler: s.mainHandler,
		},
		servo.Route{
			Path:    "/token",
			Handler: s.tokenHandler,
		},
		servo.Route{
			Path:    "/varz",
			Handler: s.statHandler,
		},
		servo.Route{
			Path:    "/cdn/",
			Handler: s.fileHandler,
		},
	}
	// setup configuration
	sopts = servo.HTTPServerOpts{
		Address:  (*ip),
		Port:     (*port),
		MaxConns: (*throtle),
		Handlers: routes,
	}
	for _, v := range routes {
		// NOTE
		// . register overwrites previously
		//   defined route iff there exists
		//   one.
		s.Register(v.Path, v.Handler)
	}
	err = s.Setup(sopts)
	if err != nil {
		panic("Server: unable to run the instance.")
	}
	fmt.Println("------------done")		
	log.Println("done setup")
	go func() {
		// wait on registered signals ( e.g. shutdown or exit )
		<-exitCh
		log.Println("* Received signal.")
		err = s.Shutdown()
		log.Printf("* Shutdown procedure is done, err=%v\n", err)
		s.MarkDone()
	}()
	log.Printf(" + started serving at (%s)", s.Address())
	err = s.Start()
	if err != nil {
		log.Println("- failed to start the server.", err)
	}
	// block on internal exit channel
	s.Wait()
	log.Println("+ exitted")
	os.Exit(0)
}
