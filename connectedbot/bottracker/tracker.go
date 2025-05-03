package bottracker

import (
	"errors"
	"io"
	"net/netip"
	"sync"
	"sync/atomic"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/connectedbot/opts"
)
var _ adapter.ConnectionTracker = (*ConnManager)(nil)

// user mean 1 config
type user struct {
	download *atomic.Int64
	upload *atomic.Int64
	Ip sync.Map
	ipCount *atomic.Int32
	maxlogin int32
	bandwidth int64
	disables *atomic.Bool
	
	allConn map[int64]ConnCloser
 	allConnAccess sync.RWMutex

	uid int //only use for alluserstatus
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("login limit cannot be zero")
)

type ConnManager struct {
	user sync.Map
	inboundManager adapter.InboundManager
	userCount *atomic.Int64
}

func NewConnManager(inmg adapter.InboundManager) *ConnManager {
	return &ConnManager{
		inboundManager: inmg,
		userCount: new(atomic.Int64),
	}
}

type ConnCloser struct {
	close io.Closer
	src netip.Addr
}

type SrInbound interface {
	DelUser(opts.ComProto) error
	AddUser(opts.ComProto) error
}

func (c *ConnManager) AddUser(u opts.User) (opts.UserStatus, error) {
	avuser, loaded := c.user.Load(u.UserStr)
	var status opts.UserStatus
	if u.MaxLogin == 0 {
		c.user.Delete(u.UserStr)
		c.removeuserinbound(u)
		return status, ErrInvalidLogin
	}
	if loaded {
		status = c.getstatus(avuser.(*user))
		return status, nil
	}
	err := c.addusertoinbound(u)
	if err != nil {
		return status, err
	}
	c.user.Store(u.UserStr, &user{
		upload: new(atomic.Int64),
		download: new(atomic.Int64),
		disables: &atomic.Bool{},
		Ip: sync.Map{},
		ipCount: new(atomic.Int32),
		maxlogin: int32(u.MaxLogin),
		bandwidth: u.Bandwidth,
		uid: u.Uid,
		allConn: map[int64]ConnCloser{},
		allConnAccess: sync.RWMutex{},

	})
	c.userCount.Add(1)
	return status, nil
}

func (c *ConnManager) GetStatusUser(u opts.User) (opts.UserStatus, error) {
	avuser, loaded := c.user.Load(u.UserStr)
	if !loaded {
		return opts.UserStatus{}, ErrUserNotFound
	}
	return c.getstatus(avuser.(*user)), nil
}

func (c *ConnManager) RemoveUser(u opts.User) (opts.UserStatus, error) { //acctualy error does not matter here it's just for to eqal with other main 4 method
	ruser, loaded := c.user.LoadAndDelete(u.UserStr)
	if !loaded {
		return opts.UserStatus{}, ErrUserNotFound
	}
	nuser := ruser.(*user)
	nuser.allConnAccess.RLock()
	for _, oconn := range nuser.allConn {
		oconn.close.Close()
	}
	nuser.allConnAccess.RUnlock()
	nuser.allConnAccess.Lock()
	nuser.allConn = map[int64]ConnCloser{}
	nuser.allConnAccess.Unlock()
	
	c.removeuserinbound(u)
	c.userCount.Add(-1)
	return c.getstatus(nuser), nil
}

func (c *ConnManager) CloseAllConn(u opts.User) {
	ruser, loaded := c.user.Load(u.UserStr)
	if !loaded {
		return
	}
	nuser := ruser.(*user)
	nuser.allConnAccess.RLock()
	for _, oconn := range nuser.allConn {
		oconn.close.Close()
	}
	nuser.allConnAccess.RUnlock()
	nuser.allConnAccess.Lock()
	nuser.allConn = map[int64]ConnCloser{}
	nuser.allConnAccess.Unlock()
}

func (c *ConnManager) ResetInbound(u opts.User) {
	c.removeuserinbound(u)
	c.addusertoinbound(u)
}

func (c *ConnManager) AddUserReset(u opts.User) (opts.UserStatus, error) {
	status, _ := c.RemoveUser(u) //error does not matter
	_, err := c.AddUser(u)
	return status, err
}

//this method do have some cost
func (c *ConnManager) AllUserStatus() (map[int]opts.UserStatus) {
	alluser := make(map[int]opts.UserStatus, c.userCount.Load())
	c.user.Range(func(key, value any) bool {
		usr := value.(*user)
		alluser[usr.uid] = c.getstatus(usr)
		return true
	})
	return alluser
}



func (c *ConnManager) addusertoinbound(u opts.User) error {

	for _, inbn := range u.InboundList {
		inbound, ok := c.inboundManager.Get(inbn)
		if !ok {
			continue
		}
		d, ok := inbound.(SrInbound)
		if ok {
			err := d.AddUser(u.Proto)
			if err != nil {
				return errors.New(u.UserStr + " user does not have any valid inbound check inbound again: " + err.Error())
			}
					
		}
	}
	return nil

	
}
func (c *ConnManager) removeuserinbound(u opts.User) {
	for _, inbound := range c.inboundManager.Inbounds() {
		d, ok := inbound.(SrInbound)
		if ok {
			d.DelUser(u.Proto)
		}
	}
}
func (c *ConnManager) getstatus(u *user) opts.UserStatus {
	if u == nil {
		return opts.UserStatus{}
	}
	status := opts.UserStatus{
		Download: u.download.Load(),
		Upload: u.upload.Load(),
		Disabled: u.disables.Load(),
		Ip: map[string]int16{},
	}
	u.Ip.Range(func(key, value any) bool {
		status.Ip[key.(netip.Addr).String()] = int16(value.(*atomic.Int64).Load())
		return true
	})
	return status
}