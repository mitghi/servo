package servo

import (
	"encoding/json"
	"errors"
)

var (
	ESRVZINVAL error = errors.New("ServerVarz: invalid/inconsistent state.")
)

// - MARK:  ServerVarz section.

func (sv *ServerVarz) ToJSON() (rep string, err error) {
	var (
		b []byte
	)
	b, err = json.Marshal(sv)
	if err != nil {
		goto ERROR
	}
	rep = string(b)

	return rep, nil
ERROR:
	return _EMPTY_, err
}

func (sv *ServerVarz) getAddr(opt byte) *int64 {
	switch opt {
	case Connection:
		return &sv.Connection
	case Active:
		return &sv.Active
	case Disconnection:
		return &sv.Disconnection
	case Total:
		return &sv.Total
	default:
		return nil
	}
}

func (sv *ServerVarz) do(inc bool, reset bool, opts byte) bool {
	var (
		fn func(*int64) = decI64
		i  byte         = 0x1
	)
	if inc {
		fn = incI64
	} else if reset {
		fn = resetI64
	}
	for ; opts > 0; i <<= 1 {
		if opts&1 == 1 {
			if addr := sv.getAddr(i); addr != nil {
				fn(addr)
			} else {
				return false
			}
		}
		opts >>= 1
	}
	return true
}

func (sv *ServerVarz) getOpts(opts byte) (p []*int64, ok bool) {
	var (
		i byte = 0x1
	)
	for ; opts > 0; i <<= 1 {
		if opts&1 == 1 {
			if addr := sv.getAddr(i); addr != nil {
				p = append(p, addr)
				ok = true
			} else {
				return nil, false
			}
		}
		opts >>= 1
	}
	return p, ok
}

func (sv *ServerVarz) Inc(opts byte) bool {
	return sv.do(true, false, opts)
}

func (sv *ServerVarz) Dec(opts byte) bool {
	return sv.do(false, false, opts)
}

func (sv *ServerVarz) Reset(opts byte) bool {
	return sv.do(false, true, opts)
}

func (sv *ServerVarz) Read(opt byte) (int64, error) {
	var (
		val *int64
	)
	val = sv.getAddr(opt)
	if val == nil {
		return 0, ESRVZINVAL
	}
	return (*val), nil
}
