package servo

type HTTPServerOpts struct {
	Address    string
	Port       int
	MaxConns   int
	WithLogger bool
	Handlers   []Route
}

func NewHTTPServerOpts() (so *HTTPServerOpts) {
	so = &HTTPServerOpts{
		Address:    "",
		Port:       0,
		MaxConns:   0,
		WithLogger: false,
	}
	return so
}
