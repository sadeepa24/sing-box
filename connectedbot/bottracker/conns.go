package bottracker

import (
	"net/netip"
	"sync"
	"sync/atomic"

	N "github.com/sagernet/sing/common/network"
)

type TCPConn struct {
	N.ExtendedConn
	connCounter *atomic.Int64
	ipMap       *sync.Map
	ipCounter   *atomic.Int32
	sourceAddr  netip.Addr
	comId       int64
	pcloser     func()
	closed      atomic.Bool
}

func (tt *TCPConn) Close() error {
	if tt.closed.Swap(true) {
		return tt.ExtendedConn.Close()
	}
	if tt.connCounter.Add(-1) <= 2 {
		tt.ipMap.Delete(tt.sourceAddr)
		tt.ipCounter.Add(-1)
	}
	if tt.pcloser != nil {
		tt.pcloser()
	}
	return tt.ExtendedConn.Close()
}

func (tt *TCPConn) Upstream() any {
	return tt.ExtendedConn
}

func (tt *TCPConn) ReaderReplaceable() bool {
	return true
}

func (tt *TCPConn) WriterReplaceable() bool {
	return true
}

type UDPConn struct {
	N.PacketConn `json:"-"`
	connCounter  *atomic.Int64
	ipMap        *sync.Map
	sourceAddr   netip.Addr
	ipCounter    *atomic.Int32
	comId        int64
	pcloser      func()
	closed       atomic.Bool
}

func (ut *UDPConn) Close() error {
	if ut.closed.Swap(true) {
		return ut.PacketConn.Close()
	}
	if ut.connCounter.Add(-1) <= 2 {
		ut.ipMap.Delete(ut.sourceAddr)
		ut.ipCounter.Add(-1)
	}
	if ut.pcloser != nil {
		ut.pcloser()
	}
	return ut.PacketConn.Close()
}

func (ut *UDPConn) Upstream() any {
	return ut.PacketConn
}

func (ut *UDPConn) ReaderReplaceable() bool {
	return true
}

func (ut *UDPConn) WriterReplaceable() bool {
	return true
}