package servo

import "testing"

func TestGetOpts(t *testing.T) {
	var (
		sv *ServerVarz = NewServerVarz()
		// cnt   *int64
		// caddr *int64 = &sv.Connection
		addrs []*int64
		ok    bool
	)
	addrs, ok = sv.getOpts(Connection)
	if !ok || len(addrs) == 0 {
		t.Fatalf("inconsistent state, unable to get addresses. addrs(%v), len(addrs(%d)), ok(%t).", addrs, len(addrs), ok)
	}
}
