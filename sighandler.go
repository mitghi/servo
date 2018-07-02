package servo

import (
	"log"
	"os"
	"sync"
)

// - MARK: Type Definition section.

/* d e b u g */
// type sigtbl map[os.Signal]func(os.Signal)

// type sigContainer struct {
// 	mu         *sync.RWMutex
// 	isCommited bool
// 	sigs       sigtbl
// 	store      []*sigHandler
// }

// type sigHandler struct {
// 	sig     os.Signal
// 	handler func(os.Signal)
// }
/* d e b u g */
// - MARK: Alloc/Init section.

func NewSigContainer() *sigContainer {
	return &sigContainer{
		mu:   &sync.RWMutex{},
		sigs: newsigtbl(),
	}
}

func NewSigHandler(sig os.Signal, handler func(os.Signal)) *sigHandler {
	return &sigHandler{
		sig:     sig,
		handler: handler,
	}
}

func newsigtbl() sigtbl {
	return sigtbl(make(map[os.Signal]func(os.Signal)))
}

// - MARK: sigtbl section.

func (stbl sigtbl) Register(sig os.Signal, fn func(os.Signal)) (isNew bool) {
	// tbl := (*stbl)
	log.Println("(Register) tbl==nil? ->", stbl == nil, stbl)
	_, isNew = stbl[sig]
	if !isNew {
		stbl[sig] = fn
		return true
	}
	return false
}

func (stbl sigtbl) Get(sig os.Signal) (fn func(os.Signal), ok bool) {
	// tbl := (*stbl)
	tbl := stbl
	fn, ok = tbl[sig]
	if !ok {
		return nil, false
	}
	return fn, true
}

func (stbl sigtbl) Remove(sig os.Signal) bool {
	var (
		ok bool
	)
	tbl := stbl
	_, ok = tbl[sig]
	if !ok {
		return false
	}
	delete(tbl, sig)
	return true
}

func (stbl sigtbl) Exists(sig os.Signal) bool {
	var (
		ok bool
	)
	tbl := stbl
	_, ok = tbl[sig]
	if !ok {
		return false
	}
	return true
}

// - MARK: SigContainer section.

func (sc *sigContainer) Add(sig os.Signal, handler func(os.Signal)) {
	// NOTE
	// . duplicate enteries in `store` container
	//   violates storage policies and cause
	//   panics.
	sc.mu.Lock()
	sc.store = append(sc.store, NewSigHandler(sig, handler))
	sc.mu.Unlock()
}

func (sc *sigContainer) Commit() (ok bool) {
	// TODO:
	// . write test case
	sc.mu.Lock()
	log.Println("sc.store->", sc.store == nil, sc.store)
	for _, v := range sc.store {
		ok = sc.sigs.Register(v.sig, v.handler)
		if !ok {
			sc.mu.Unlock()
			goto ERROR
		}
	}
	sc.mu.Unlock()
	return true
ERROR:
	return false
}
