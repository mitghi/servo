package servo

// - DEBUG STASH - //

/* d e b u g */
// type NewConnectionHandler func(conn net.Conn) interface{}
// type ClientConnectionHandler func(conn net.Conn, client interface{})
// type ClientDisconnectionHandler func(conn net.Conn, client interface{})
// type ServerStartCallback func(s *Server, sid string)
// type ServerStopCallback func(s *Server, sid string)
// type StatsFunc func(net.Conn)

// type clientHandler struct {
// 	onNewConnection NewConnectionHandler

// 	onServerStart         ServerStartCallback
// 	onServerStop          ServerStopCallback
// 	onClientConnection    ClientConnectionHandler
// 	onClientDisconnection ClientDisconnectionHandler
// }

// // ServerVarz is a 64bit aligned atomic struct
// // for statistics. It is used by `Server` instance
// // to keep track of various variables.
// type ServerVarz struct {
// 	Connection    int64 `json:"connection"`
// 	Active        int64 `json:"active"`
// 	Disconnection int64 `json:"disconnection"`
// 	Total         int64 `json:"total"`
// }

// // ServerStats is the struct that implements
// // statistics handling. It serves statistics
// // over HTTP.
// type ServerStats struct {
// 	conn net.Conn
// 	mu   *sync.RWMutex
// 	varz *ServerVarz
// }
/* d e b u g */

/* d e b u g */
// // ServerStatsInterface is the interface that
// // must be conformed by implementors delgating
// // statistics procs.
// type ServerStatsInterface interface {
// 	GetStats() *ServerVarz
// 	Spawn() bool
// }

// // Server is the main struct that implements
// // serving functionalities. It can be configured
// // in various modes with different serving behaviors.
// type Server struct {
// 	status uint32

// 	clientHandler
// 	listener  net.Listener
// 	mu        sync.Mutex
// 	wgMu      sync.Mutex
// 	apDoneCh  chan struct{}
// 	wg        sync.WaitGroup
// 	optsTLS   *TLSOpts
// 	Varz      *ServerVarz
// 	SID       string
// 	Addr      string
// 	host      string
// 	port      int
// 	startDL   int
// 	doneSetup bool
// 	doneInit  bool
// 	useTLS    bool
// }

// // ServerOpts is a container for holding configuration.
// // It is passed to the server for initial setup.
// type ServerOpts struct {
// 	OnServerStart         ServerStartCallback
// 	OnServerStop          ServerStopCallback
// 	OnNewConnection       NewConnectionHandler
// 	OnClientConnection    ClientConnectionHandler
// 	OnClientDisconnection ClientDisconnectionHandler
// 	TLSOptions            *TLSOpts
// 	Addr                  string
// 	StartDeadline         int
// 	UseTLS                bool
// }
/* d e b u g */

/* d e b u g */
// func TestPLG(t *testing.T) {
// 	var (
// 		f   *os.File
// 		err error
// 	)
// 	f, err = os.OpenFile("/tmp/test.txt", os.O_RDWR, 0600)
// 	if err != nil {
// 		f, err = os.Create("/tmp/test.txt")
// 		if err != nil {
// 			t.Fatal("unable to open the file.", err)
// 		}
// 		t.Fatal("unable to open the file.", err)

// 	}
// 	log.SetOutput(f)
// 	log.Println("sample output")
// }
/* d e b u g */
