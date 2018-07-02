package servo

import (
	"errors"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mitghi/util/atoms"
)

/**
* TODO:
* . profile mutex
* . write test case for mutex locks ( i.e. check for deadlocks )
* . implement Stat/Callback mechanism
**/

// UNIMPLEMENTED:
// . type StatFn func(instance *HTTPServer, supervisor *Server, evname string, data interface{})

// ServerMode is the serving mode identifier
type ServerMode byte

// Serving modes
const (
	ModeTCP ServerMode = iota
	ModeTLS
)

// Error messages
var (
	EINVALAddr  error = errors.New("server: invalid address")
	EINVAL      error = errors.New("server: invalid")
	EFATAL      error = errors.New("server: encountered fatal condition")
	EINCOMP     error = errors.New("server: incomplete")
	EFATALSTART error = errors.New("server: fatal failure in server start")
)

// Error messages (Server instance)
var (
	ESRVFATALOPS error = errors.New("server: fatal error during operation")
)

// Server status codes
const (
	SRVNONE uint32 = iota
	SRVRUNNING
	SRVPAUSE
	SRVGODOWN
	SRVFATAL
)

// Program defaults
const (
	_EMPTY_           string = ""
	_FATAL_           string = "FATAL"
	cTLPROTO          string = "tcp"
	cDefaultLogPath   string = "/tmp/servo.out"
	cDefSTARTDEADLINE int    = 2
	cMINNEGRACE              = time.Millisecond * 20
	cMAXNEGRACE              = time.Second * 2
)

// Connection-Status bit flags
const (
	Connection    = 1 << iota // 1
	Active        = 1 << iota // 2
	Disconnection = 1 << iota // 4
	Total         = 1 << iota // 8
)

// Route associates the path from
// a given URL to a particular `HandlerFn`.
// It is used for stashing routes before
// commiting them to a hash table (
// for maintanibility purposes ).
type Route struct {
	Path    string
	Handler HandlerFn
}

// TODO:
// . protocol conformance
type HTTPServerInterface interface {
	Setup(HTTPServerOpts) error
	Shutdown() error
	Start() error
	SubmitRoutes() bool
	Wait()
}

// HandlerFn is the type identifier for functions
// used in handling incomming requests.
type HandlerFn func(w http.ResponseWriter, r *http.Request)

// HTTPServer is a struct that wraps around
// `Server` and serves HTTP requests.
type HTTPServer struct {
	// size: 128 bytes ( aligned )
	// _pad0      int64 // padding
	// _pad1      int64 // padding
	_pad       [2]int64
	DoneCh     chan struct{}
	exitCh     chan os.Signal
	Varz       *ServerVarz
	hasSignals *atoms.Boolean
	running    *atoms.Boolean
	doneInit   *atoms.Boolean
	mu         *sync.RWMutex
	server     *http.Server
	mux        *http.ServeMux
	parent     *Server
	sigbox     *sigContainer
	routes     *Routes
	Logger     *LogFile
	addr       string
	// TODO:
	// . implement callbacks
	// . . callbacks  map[string]StatFn
}

// Delegate types
type (
	// NewConnectionHandler is the delegate function responsible
	// of allocating and initializing incomming connections.
	NewConnectionHandler func(net.Conn) interface{}
	// ClientConnectionHandler is the delegate function
	// which gets invoked AFTER initialization of
	// a new connection.
	ClientConnectionHandler func(net.Conn, interface{})
	// ClientDisconnectionHandler is the delegate function
	// which gets invoked when client connection terminates.
	ClientDisconnectionHandler func(net.Conn, interface{})
	// ServerStartCallback is the delegate function
	// which gets invoked when Server start procedure
	// is invoked.
	ServerStartCallback func(*Server, string)
	// ServerStopCallback is the delegate function
	// which gets invoked when Server stop procedure
	// is invoked.
	ServerStopCallback func(*Server, string)
	// StatsFunc is the delegate type
	// which provides statistical infos.
	StatsFunc func(net.Conn)
)

// clientHandler is the container which
// holds delegate functions for handling
// underlaying `Client` implementation. It
// is stored as a node and meant to be
// embedded by a parent struct.
type clientHandler struct {
	onNewConnection NewConnectionHandler

	onServerStart         ServerStartCallback
	onServerStop          ServerStopCallback
	onClientConnection    ClientConnectionHandler
	onClientDisconnection ClientDisconnectionHandler
}

// ServerVarz is a 64bit aligned atomic struct
// for statistics. It is used by `Server` instance
// to keep track of various variables.
type ServerVarz struct {
	Connection    int64 `json:"connection"`
	Active        int64 `json:"active"`
	Disconnection int64 `json:"disconnection"`
	Total         int64 `json:"total"`
}

// ServerStats is the struct that implements
// statistics handling. It serves statistics
// over HTTP.
type ServerStats struct {
	conn net.Conn
	mu   *sync.RWMutex
	varz *ServerVarz
}

// ServerStatsInterface is the interface that
// must be conformed by implementors delgating
// statistics procs.
type ServerStatsInterface interface {
	GetStats() *ServerVarz
	Spawn() bool
}

// Server is the main struct that implements
// serving functionalities. It can be configured
// in various modes with different serving behaviors.
// TODO:
// . struct alignment
type Server struct {
	status uint32

	clientHandler
	listener  net.Listener
	mu        sync.Mutex
	wgMu      sync.Mutex
	apDoneCh  chan struct{}
	wg        sync.WaitGroup
	optsTLS   *TLSOpts
	Varz      *ServerVarz
	SID       string
	Addr      string
	host      string
	port      int
	startDL   int
	doneSetup bool
	doneInit  bool
	useTLS    bool
}

// ServerOpts is a container for holding configuration.
// It is passed to the server for initial setup.
type ServerOpts struct {
	OnServerStart         ServerStartCallback
	OnServerStop          ServerStopCallback
	OnNewConnection       NewConnectionHandler
	OnClientConnection    ClientConnectionHandler
	OnClientDisconnection ClientDisconnectionHandler
	TLSOptions            *TLSOpts
	Addr                  string
	StartDeadline         int
	UseTLS                bool
}

// - MARK: Signal section.

// sigtbl is the type for Signal Table.
type sigtbl map[os.Signal]func(os.Signal)

// sigContainer is the struct that implements
// signal registration and acts as a stash
// container for signals.
type sigContainer struct {
	mu         *sync.RWMutex
	isCommited bool
	sigs       sigtbl
	store      []*sigHandler
}

// sigHandler is the struct that
// associates a particular `os.Signal`
// to a `func(os.Signal)` function.
type sigHandler struct {
	sig     os.Signal
	handler func(os.Signal)
}
