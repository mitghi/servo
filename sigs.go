package servo

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// - MARK: Signal-Handlers section.`

func (s *HTTPServer) RegisterSignalHandlers() {
	var (
		ok bool
	)
	s.mu.Lock()
	for _, v := range s.sigbox.store {
		ok = s.sigbox.sigs.Register(v.sig, v.handler)
		if !ok {
			log.Println("(RegisterSignalHandlers) Duplicate table entry / Violation.")
			return
		}
	}
	s.mu.Unlock()
}

func (s *HTTPServer) RegisterSignalHandler(sig os.Signal, fn func(os.Signal)) {
	// NOTE:
	// . handlers contained by `store` must be commited
	//   prior to running the server instance, malfunction
	//   should be expected otherwise.
	s.sigbox.Add(sig, fn)
}

func (s *HTTPServer) handleSIGHUP(sig os.Signal) {
	s.mu.Lock()
	s.Varz.Reset(Connection | Active | Disconnection | Total)
	s.mu.Unlock()
}

func (s *HTTPServer) handleSignals() {
	var (
		sig os.Signal
		err error
	)
	sig = <-s.exitCh
	err = s.Shutdown()
	// TODO:
	// . refactor into a better implementation ( use discarded value )
	_ = s.dispatchSig(sig)
	log.Printf("* Shutdown procedure is done, err=%+v. \n", err)
	s.MarkDone()
}

func (s *HTTPServer) HandleSignals() bool {
	if s.hasSignals.Is(true) {
		return false
	}
	signal.Notify(s.exitCh, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go s.handleSignals()
	s.hasSignals.Set(true)
	return true
}

func (s *HTTPServer) dispatchSig(sig os.Signal) (err error) {
	var (
		handler func(os.Signal)
		ok      bool
	)
	handler, ok = s.sigbox.sigs.Get(sig)
	if handler == nil || !ok {
		log.Printf("(dispatchSig) No handler is registered for %+v sig\n.", sig)
		err = ESRVINVKFATAL
		goto ERROR
	}
	/* d e b u g */
	// switch sig {
	// case syscall.SIGHUP:
	// 	handler, ok = s.sigbox.sigs.Get(sig)
	// 	if handler == nil || !ok {
	// 		log.Printf("(dispatchSig) No handler is registered for %+v sig\n.", sig)
	// 		err = ESRVINVKFATAL
	// 		goto ERROR
	// 	}
	// 	break
	// default:
	// 	err = ESRVLKUPFATAL
	// 	goto ERROR
	// }
	/* d e b u g */
	// invoke signal handler
	handler(sig)
	return nil
ERROR:
	return err
}
