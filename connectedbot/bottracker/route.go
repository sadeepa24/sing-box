package bottracker

import (
	"context"
	"math/rand"
	"net"
	"sync/atomic"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing/common/bufio"
	N "github.com/sagernet/sing/common/network"
)

//TODO: add another function for comon things in both RoutedConnection & RoutedPacketConnection
func (c *ConnManager) RoutedConnection(_ context.Context, conn net.Conn, metadata adapter.InboundContext, _ adapter.Rule, _ adapter.Outbound) net.Conn {
	ruser, loaded := c.user.Load(metadata.User)
	if !loaded {
		conn.Close()
		return conn
	}
	muser := ruser.(*user)
	if muser.disables.Load() {
		conn.Close()
		return conn
	}

	if muser.download.Load() + muser.upload.Load() >= muser.bandwidth {
		if muser.disables.Swap(true) {
			//recheck 
			conn.Close()
			return conn
		}
		muser.allConnAccess.Lock()
		for _, oconn := range muser.allConn {
			oconn.close.Close()
		}
		muser.allConn = map[int64]ConnCloser{}
		muser.allConnAccess.Unlock()
		conn.Close()


		return conn
	}
	

	connCount, ipvalid := muser.Ip.Load(metadata.Source.Addr)
	if !ipvalid {
		if muser.ipCount.Load() >= muser.maxlogin {
			conn.Close()
			return conn
		}
		connCount = new(atomic.Int64)
		muser.Ip.Store(metadata.Source.Addr, connCount)
		muser.ipCount.Add(1)
	}

	connCount.(*atomic.Int64).Add(1)

	id := rand.Int63()
	nconn := &TCPConn{
		ExtendedConn: bufio.NewCounterConn(conn, []N.CountFunc{func(n int64) {
			muser.upload.Add(n)
		}}, []N.CountFunc{func(n int64) {
			if muser.download.Add(n) > muser.bandwidth {
				conn.Close()
			}
		}}),
		connCounter: connCount.(*atomic.Int64),
		sourceAddr: metadata.Source.Addr,
		ipMap: &muser.Ip,
		comId: id,
		ipCounter: muser.ipCount,
		pcloser: func ()  {

			if muser.download.Load() + muser.upload.Load() >= muser.bandwidth {
				if muser.disables.Swap(true) {
					return
				}
				muser.allConnAccess.Lock()
				if len(muser.allConn) == 0 {
					muser.allConnAccess.Unlock()
					return
				}
				for _, oconn := range muser.allConn {
					oconn.close.Close()
				}
				muser.allConn = map[int64]ConnCloser{}
				muser.allConnAccess.Unlock()
				return
			}

			muser.allConnAccess.RLock()
			_, exists := muser.allConn[id]
			muser.allConnAccess.RUnlock()
			if exists {
				muser.allConnAccess.Lock()
				delete(muser.allConn, id)
				muser.allConnAccess.Unlock()
			}

		},
	} 

	muser.allConnAccess.Lock()
	muser.allConn[nconn.comId] = ConnCloser{close: conn,src: metadata.Source.Addr,}
	muser.allConnAccess.Unlock()

	return nconn
}


func (c *ConnManager) RoutedPacketConnection(_ context.Context, conn N.PacketConn, metadata adapter.InboundContext, _ adapter.Rule, _ adapter.Outbound) N.PacketConn {
	

	ruser, loaded := c.user.Load(metadata.User)
	if !loaded {
		conn.Close()
		return conn
	}
	muser := ruser.(*user)
	if muser.disables.Load() {
		conn.Close()
		return conn
	}

	if muser.download.Load() + muser.upload.Load() >= muser.bandwidth {
		if muser.disables.Swap(true) {
			//recheck 
			conn.Close()
			return conn
		}
		muser.allConnAccess.Lock()
		for _, oconn := range muser.allConn {
			oconn.close.Close()
		}
		muser.allConn = map[int64]ConnCloser{}
		muser.allConnAccess.Unlock()
		conn.Close()
		return conn
	}
	

	connCount, ipvalid := muser.Ip.Load(metadata.Source.Addr)
	if !ipvalid {
		if muser.ipCount.Load() >= muser.maxlogin {
			conn.Close()
			return conn
		}
		connCount = new(atomic.Int64)
		muser.Ip.Store(metadata.Source.Addr, connCount)
		muser.ipCount.Add(1)
	}

	connCount.(*atomic.Int64).Add(1)

	id := rand.Int63()
	
	nconn := &UDPConn{
		PacketConn: bufio.NewCounterPacketConn(conn, []N.CountFunc{func(n int64) {
			muser.upload.Add(n)
		}}, []N.CountFunc{func(n int64) {
			if muser.download.Add(n) > muser.bandwidth {
				conn.Close()
			}
		}}),
		connCounter: connCount.(*atomic.Int64),
		sourceAddr: metadata.Source.Addr,
		ipMap: &muser.Ip,
		comId: id,
		ipCounter: muser.ipCount,
		pcloser: func ()  {

			if muser.download.Load() + muser.upload.Load() >= muser.bandwidth {
				if muser.disables.Swap(true) {
					return
				}
				muser.allConnAccess.Lock()
				if len(muser.allConn) == 0 {
					muser.allConnAccess.Unlock()
					return
				}
				for _, oconn := range muser.allConn {
					oconn.close.Close()
				}
				muser.allConn = map[int64]ConnCloser{}
				muser.allConnAccess.Unlock()
				return
			}

			muser.allConnAccess.RLock()
			_, exists := muser.allConn[id]
			muser.allConnAccess.RUnlock()
			if exists {
				muser.allConnAccess.Lock()
				delete(muser.allConn, id)
				muser.allConnAccess.Unlock()
			}

		},
	}
	muser.allConnAccess.Lock()
	muser.allConn[nconn.comId] = ConnCloser{close: conn,src: metadata.Source.Addr,}
	muser.allConnAccess.Unlock()

	return nconn

}