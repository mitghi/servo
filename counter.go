package servo

import (
	"log"
	"sync"
)

/**
* TODO:
* . implement generic counter
**/

// - MARK: Counter section.

// CounterEvent is the type for Counter
type CounterEvent byte

// Events enumeration for Counter
const (
	CEAdded CounterEvent = iota
	CERemoved
	CEFetched
)

// CFn is the type identifier for Counter callback function
type CFn func(string, map[string]int)

// Counter implements a Expiring Set which triggers
// an event when the threshold is reached.
type Counter struct {
	mu        *sync.RWMutex
	c         map[string]int
	threshold int
	callback  CFn
}

// NewCounter allocates and initializes a new `Counter`
// struct and returns a new pointer to it.
func NewCounter() *Counter {
	return &Counter{
		mu: &sync.RWMutex{},
		c:  make(map[string]int),
	}
}

// Add increments count associated to `key`. It
// triggers the callback when `count>threshold`.
func (c *Counter) Add(key string) {
	var (
		count int
		ok    bool
	)
	log.Println("(counter) add is invoked. callback ->", c.callback)
	c.mu.Lock()
	count, ok = c.c[key]
	if !ok {
		c.c[key] = 1
	} else {
		count += 1
		c.c[key] = count
		if count > c.threshold && c.callback != nil {
			c.callback(key, c.c)
			c.Log(CERemoved)
		}
	}
	c.mu.Unlock()
}

// Removes removes associated counter to `key`
// from the set.
func (c *Counter) Remove(key string) (ok bool) {
	c.mu.Lock()
	_, ok = c.c[key]
	if ok {
		delete(c.c, key)
	}
	c.mu.Unlock()
	return ok
}

// Get gets associated counter value to `key`
// from the set along with a boolean to indicate
// success status.
func (c *Counter) Get(key string) (value int, ok bool) {
	c.mu.Lock()
	value, ok = c.c[key]
	c.mu.Unlock()
	return value, ok
}

// Clear purges all enteries in the set.
func (c *Counter) Clear() {
	c.mu.Lock()
	for k, _ := range c.c {
		delete(c.c, k)
	}
	c.mu.Unlock()
}

// Exists checks whether `key` is present inside the
// set.
func (c *Counter) Exists(key string) (ok bool) {
	c.mu.RLock()
	_, ok = c.c[key]
	c.mu.RUnlock()
	return ok
}

// SetCallback registers `fn` as callback function.
// It gets invoked when a k/v pair passes `threshold`
// value.
func (c *Counter) SetCallback(fn CFn) {
	c.mu.Lock()
	c.callback = fn
	c.mu.Unlock()
}

// SetThreshold sets the threshold which acts
// as maximum counting limit before trigerring
// the callback function.
func (c *Counter) SetThreshold(threshold int) {
	c.mu.Lock()
	c.threshold = threshold
	c.mu.Unlock()
}

func (c *Counter) Log(ev CounterEvent) {
	// TODO:
	switch ev {
	case CEAdded:
		log.Println("(Counter) Key is added.")
		// TODO:
		// . append that to a log file.
	case CERemoved:
		log.Println("(Counter) Key is removed.")
		// TODO:
		// . invoke callback
	case CEFetched:
	default:
		log.Println("(Counter) undefined Counter Event. ")
		return
	}
}

// - MARK: EventCounter section.

type EventCounter struct {
	Events []interface{}
	Count  int
}

func NewEventCounter() (ec *EventCounter) {
	ec = &EventCounter{
		Events: make([]interface{}, 0),
		Count:  0,
	}
	return ec
}

func (ec *EventCounter) Push(item interface{}) (ok bool) {
	ec.Events = append(ec.Events, item)
	ec.Count += 1
	// TODO:
	// . return correct status indicator
	//   when slice memory allocation
	//   fails.
	return true
}

func (ec *EventCounter) Pop() (item interface{}) {
	if ec.Count == 0 {
		return nil
	}
	item = ec.Events[0]
	copy(ec.Events, ec.Events[1:])
	ec.Count -= 1
	return item
}

func (ec *EventCounter) Size() int {
	return ec.Count
}

// Head returns head of the list without removing
// it. It returns a `null` pointer when no item
// exists.
func (ec *EventCounter) Head() (item interface{}) {
	if ec.Count == 0 {
		return nil
	}
	item = ec.Events[0]
	return item
}

func (ec *EventCounter) Tail() (item interface{}) {
	if ec.Count == 0 {
		return nil
	}
	item = ec.Events[ec.Count-1]
	return item
}
