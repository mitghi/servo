// servobase provides protocol (interface) definitions for
// interpackage compatibility.
package servobase

// ServerStatsInterface is the interface that
// must be conformed by implementors delegating
// statistics procs.
type ServerStatsInterface interface {
	GetStats() *ServerVarz
	Spawn() bool
}

// HTTPServerInterface is the interface
// that the implementor must conform to
// if delegation is desired.
type HTTPServerInterface interface {
	Setup(HTTPServerOpts) error
	Shutdown() error
	Start() error
	SubmitRoutes() bool
	MarkDone()
	Wait()
}
