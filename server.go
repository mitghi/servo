package servo

import (
	"log"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

/**
* TODO:
* . finish documentation
**/

// - MARK: Counter section.

func incI64(addr *int64) {
	atomic.AddInt64(addr, 1)
}

func decI64(addr *int64) {
	atomic.AddInt64(addr, -1)
}

func resetI64(addr *int64) {
	atomic.StoreInt64(addr, 0)
}

// - MARK: Type definition section.

func NewServerStats() *ServerStats {
	var (
		s *ServerStats
	)
	s = &ServerStats{
		mu:   &sync.RWMutex{}, // read/write mutex to protect parallel access
		conn: nil,             // net.Conn underlaying tcp connection
		varz: &ServerVarz{
			Connection:    0,
			Active:        0,
			Disconnection: 0,
			Total:         0,
		},
	}
	return s
}

// GetStatus is a receiver function that returns
// relevant statistic variable. Note, the returned
// value is a copy.
func (s *ServerStats) GetStatus() (*ServerVarz, error) {
	var (
		sv *ServerVarz = NewServerVarz()
	)
	s.mu.Lock()
	if s.varz == nil {
		s.mu.Unlock()
		return nil, ESRVFATALOPS
	}
	sv.Active = atomic.LoadInt64(&s.varz.Active)
	sv.Connection = atomic.LoadInt64(&s.varz.Connection)
	sv.Disconnection = atomic.LoadInt64(&s.varz.Disconnection)
	sv.Total = atomic.LoadInt64(&s.varz.Total)
	s.mu.Unlock()
	return sv, nil
}

// NewServerVarz allocates and initializes a new
// `ServerVarz` struct and returns a pointer to
// it. Note, this struct is 64bit aligned for
// using atomics.
func NewServerVarz() (sv *ServerVarz) {
	sv = &ServerVarz{
		Connection:    0,
		Active:        0,
		Disconnection: 0,
		Total:         0,
	}
	return sv
}

// Register is a receiver function that registers
// ahthe given `path` with `handler`. It registers
// new routes to the underlaying server framework.
func (s *ServerStats) Register(path string, handler StatsFunc) bool {
	// TODO:
	return false
}

// Spawn activates a new server instance and returns
// a boolean to indicate its success.
func (s *ServerStats) Spawn() bool {
	// TODO:
	return false
}

// - MARK: alloc/init section.

// NewServer is a function that allocates and initializes
// a new server instance and returns a pointer to it.
func NewServer() (s *Server) {
	var (
		sid string
		err error
	)
	s = &Server{
		status:    SRVNONE,
		Varz:      &ServerVarz{},
		doneSetup: false,
		doneInit:  false,
	}
	sid, err = RandString16()
	if err == nil {
		goto OK
	}
	sid, err = RandString(16)
	if err == nil {
		goto OK
	}
	s.SID = _FATAL_

	return s
OK:
	s.SID = sid

	return s
}

func NewServerFromConfig(opts ServerOpts) (s *Server, err error) {
	s = NewServer()
	if err = s.ApplyConfig(opts); err != nil {
		return nil, err
	}

	return s, nil
}

// - MARK: Server section.

// setServerOpts sets and validates `ServerOpts`
// and returns a boolean to indicate its validity.
// when `s` is given, it sets the associated variables
// on the `Server` instance accordingly.
func setServerOpts(s *Server, opts ServerOpts) (ok bool, err error) {
	var (
		h            string                     // host
		ps           string                     // port
		p            int                        // ps as integer
		startDL      int    = cDefSTARTDEADLINE // server start deadline
		shouldCommit bool   = s != nil          // mutate `s` when true
	)
	if len(opts.Addr) == 0 {
		err = EINCOMP
		goto ERROR
	}
	h, ps, err = net.SplitHostPort(opts.Addr)
	if err != nil {
		goto ERROR
	}
	if len(ps) == 0 {
		err = EINCOMP
		goto ERROR
	}
	p, err = strconv.Atoi(ps)
	if err != nil {
		goto ERROR
	}
	if opts.StartDeadline > 0 {
		startDL = opts.StartDeadline
	}
	if opts.UseTLS && opts.TLSOptions == nil {
		err = EINCOMP
		goto ERROR
	}
	if shouldCommit {
		if s == nil {
			err = EINVAL
			goto ERROR
		}
		s.clientHandler.onNewConnection = opts.OnNewConnection
		s.clientHandler.onClientConnection = opts.OnClientConnection
		s.clientHandler.onClientDisconnection = opts.OnClientDisconnection
		s.clientHandler.onServerStart = opts.OnServerStart
		s.clientHandler.onServerStop = opts.OnServerStop
		s.optsTLS = opts.TLSOptions
		s.startDL = startDL
		s.Addr = opts.Addr
		s.host = h
		s.port = p
		s.useTLS = opts.UseTLS
	}

	return true, nil
ERROR:
	return false, err
}

// ApplyConfig is a receiver method that performs
// initial configuration parsing from the given
// server options `opts`. It returns an error to
// indicate failure.
func (s *Server) ApplyConfig(opts ServerOpts) (err error) {
	if s.doneSetup {
		return EFATAL
	}
	if ok, err := setServerOpts(s, opts); (!ok) || (err != nil) {
		return err
	}
	s.doneSetup = true

	return nil
}

func (s *Server) initializeServer() {
	// TODO:
	s.apDoneCh = make(chan struct{}, 1)
}

func (s *Server) SetStatus(status uint32) {
	atomic.StoreUint32(&s.status, status)
}

func (s *Server) GetStatus() uint32 {
	return atomic.LoadUint32(&s.status)
}

func (s *Server) IsRunning() bool {
	return s.GetStatus() == SRVRUNNING
}

func (s *Server) handleNewConnection(conn net.Conn) interface{} {
	log.Println("[connectionHandler]: new connection established")
	return nil
}

func (s *Server) SetRunning(r bool) {
	if r {
		s.SetStatus(SRVRUNNING)
	} else {
		s.SetStatus(SRVNONE)
	}
}

func (s *Server) errorHandler(status uint32, caller string) {
	// TODO:
	// . replace caller with binary flags
}

func (s *Server) handleLoop(readyCh chan<- struct{}) (err error) {
	defer func() {
		if err != nil && readyCh != nil {
			close(readyCh)
		}
	}()
	if s.GetStatus() != SRVNONE {
		err = EINVAL
		return err
	}
	var (
		sid         string               // server identifier
		status      uint32               // server status
		l           net.Listener         // listener fd
		connHandler NewConnectionHandler // connection handler
		neGrace     time.Duration        // network error grace period
		doneCh      chan struct{}        // chan to signal end of procedure
	)
	s.mu.Lock()
	l = s.listener
	doneCh = s.apDoneCh
	if connHandler = s.clientHandler.onNewConnection; connHandler == nil {
		connHandler = s.handleNewConnection
	}
	sid = s.SID
	s.mu.Unlock()
	if l == nil || doneCh == nil || connHandler == nil {
		s.SetStatus(SRVFATAL)
		return EFATAL
	}
	s.SetStatus(SRVRUNNING)
	// signal loop start
	readyCh <- struct{}{}
	readyCh = nil
	// call delegate method for
	// server start.
	if s.clientHandler.onServerStart != nil {
		s.clientHandler.onServerStart(s, sid)
	}
ML:
	for s.IsRunning() {
		var (
			c   net.Conn  // remote connection
			eok bool      // boolean used in typecasting error messages
			ne  net.Error // network error
		)
		c, err = l.Accept()
		if err != nil {
			if ne, eok = err.(net.Error); eok && ne.Temporary() {
				time.Sleep(neGrace)
				neGrace *= 2
				if neGrace > cMAXNEGRACE {
					neGrace = cMAXNEGRACE
				}
			} else if s.IsRunning() {
				log.Println("[server]: discarding temporary connection error")
			} else {
				status = s.GetStatus()
				switch status {
				case SRVFATAL, SRVGODOWN, SRVNONE:
					log.Println("[server]: fatal-error/down-flag in server loop")
					break ML
				default:
				}
			}
			continue
		}
		neGrace = cMINNEGRACE
		s.wgMu.Lock()
		s.wg.Add(1)
		go func() {
			connHandler(c)
			s.wg.Done()
		}()
		s.wgMu.Unlock()
	}
	if status == SRVNONE || status == SRVRUNNING {
		status = s.GetStatus()
	}
	switch status {
	case SRVFATAL:
		log.Println("[server]: status is (SRVFATAL) before shutdown")
	default:
	}
	if doneCh != nil {
		doneCh <- struct{}{}
	}

	return err
}

func (s *Server) Start() (err error) {
	var (
		l         net.Listener  // listening fd
		readyChan chan struct{} // used by handleLoop(....) to signal readiness
		done      bool          // prevent inconsistency when timer goes off
	)
	if s.GetStatus() != SRVNONE {
		err = EFATAL
		goto ERROR
	}
	s.mu.Lock()
	if !s.doneInit {
		s.initializeServer()
		l, err = net.Listen(cTLPROTO, s.Addr)
		if err != nil {
			s.cleanUp()
			s.SetStatus(SRVFATAL)
			s.mu.Unlock()
			return err
		}
		s.listener = l
		s.doneInit = true
	}
	s.mu.Unlock()
	readyChan = make(chan struct{})
	go s.handleLoop(readyChan)
	select {
	case _ = <-time.After(time.Second * time.Duration(s.startDL)):
		if !done {
			err = EFATAL
			goto ERROR
		}
	case _, ok := <-readyChan:
		done = true
		if !ok {
			err = EFATALSTART
		}
		break
	}

	return nil
ERROR:
	return err
}

func (s *Server) Stop() (ok bool, err error) {
	var (
		sid      string        // server identifier
		listener net.Listener  // underlaying socket fd
		doneCh   chan struct{} // signals handleloop termination
		cok      bool
	)
	/* d e b u g */
	// defer func() {
	// 	log.Println("[server(stop)]: error before exit(", err, ")")
	// }()
	/* d e b u g */
	if s.GetStatus() != SRVRUNNING {
		err = EINVAL
		goto ERROR
	}
	s.mu.Lock()
	listener = s.listener
	doneCh = s.apDoneCh
	sid = s.SID
	s.mu.Unlock()
	if listener == nil {
		err = EFATAL
		goto ERROR
	}
	s.listener = nil
	if doneCh == nil {
		err = EFATAL
		goto ERROR
	}
	s.SetStatus(SRVGODOWN)
	_ = listener.Close()
	_, cok = <-doneCh
	if cok {
		s.mu.Lock()
		s.cleanUp()
		s.mu.Unlock()
	} else {
		log.Println("[server]: handloop is not active")
	}
	if s.clientHandler.onServerStop != nil {
		s.clientHandler.onServerStop(s, sid)
	}
	s.wgMu.Lock()
	// wait for coroutines to finish
	s.wg.Wait()
	s.wgMu.Unlock()

	return true, nil
ERROR:
	return false, err
}

func (s *Server) cleanUp() {
	if s.apDoneCh != nil {
		close(s.apDoneCh)
		s.apDoneCh = nil
	}
	s.doneInit = false
}
