package servo

import (
	"log"
	"net"
	"testing"
	"time"
)

func TestServerInit(t *testing.T) {
	var (
		s      *Server
		connfn func(net.Conn) interface{}
		opts   ServerOpts
		ok     bool
		err    error
	)
	connfn = func(conn net.Conn) interface{} {
		log.Println("in dummy function")
		return nil
	}
	opts = ServerOpts{
		OnNewConnection: connfn,
		Addr:            ":8080",
		StartDeadline:   2,
	}
	if ok, err := setServerOpts(nil, opts); !ok || err != nil {
		t.Fatal("inconsistent state, expected err==nil.")
	}
	s, err = NewServerFromConfig(opts)
	if err != nil {
		t.Fatal("inconsistent state, expected err==nil.")
	}
	if s.clientHandler.onNewConnection == nil {
		t.Fatal("assertion failed, expected non-null value.")
	}
	_ = s.clientHandler.onNewConnection(nil)
	err = s.Start()
	if err != nil {
		t.Fatal("inconsistent state, expected null. cannot start server.", err)
	}
	time.Sleep(4)
	ok, err = s.Stop()
	if !ok || err != nil {
		t.Fatal("assertion failed, expected true and null.")
	}
}

func TestServerVarz(t *testing.T) {
	var (
		sv  *ServerVarz = &ServerVarz{}
		err error
	)
	sv.Inc(Connection | Total)
	if (sv.Connection == sv.Total && sv.Connection != 1) || (sv.Connection != sv.Total) {
		t.Fatal("inconsistent state, expected equal.")
	}
	if (sv.Disconnection == sv.Active && sv.Active != 0) || (sv.Disconnection != sv.Active) {
		t.Fatal("inconsistent state, expected equal.")
	}
	sv.Inc(Connection)
	sv.Inc(Connection)
	sv.Inc(Connection)
	sv.Dec(Total)
	if sv.Connection != 4 && sv.Total != 0 {
		t.Fatal("inconsistent state, expected equal.")
	}
	str, err := sv.ToJSON()
	if err != nil {
		t.Fatal("inconsistent state, expected err==nil.", err)
	}
	log.Println("json repr: ", str)
}
