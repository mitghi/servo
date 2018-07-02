package servo

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Routes struct {
	mu     *sync.RWMutex
	routes map[string]HandlerFn
}

// - MARK: Routes section.

func (r *Routes) Register(path string, handler HandlerFn) (ok bool) {
	r.mu.Lock()
	_, ok = r.routes[path]
	if !ok {
		r.routes[path] = handler
		ok = true
	}
	r.mu.Unlock()
	return ok
}

func (r *Routes) Exists(path string) (ok bool) {
	r.mu.RLock()
	_, ok = r.routes[path]
	r.mu.RUnlock()
	return ok
}

func (r *Routes) Get(path string) (handler HandlerFn, ok bool) {
	r.mu.RLock()
	handler, ok = r.routes[path]
	r.mu.RUnlock()
	if ok {
		return handler, true
	}
	return nil, false
}

func (r *Routes) Remove(path string) (ok bool) {
	r.mu.Lock()
	_, ok = r.routes[path]
	if ok {
		delete(r.routes, path)
	}
	r.mu.Unlock()
	return ok
}

// SubmitToMux is a receiver method that commits
// submited routes into the underlaying '*http.ServeMux'
// instance. Note, it panics when the supplied '*http.ServeMux'
// instance is nil.
func (r *Routes) SubmitToMux(mux *http.ServeMux) bool {
	log.Println("Routes: length of r.routes", len(r.routes), mux == nil)
	if mux == nil {
		panic("Routes: http.mux==nil")
	}
	fmt.Println("routes", r.routes)
	for path, handler := range r.routes {
		mux.HandleFunc(path, handler)
		/* d e b u g */
		// log.Println("submited", path, handler)
		/* d e b u g */
	}
	return true
}
