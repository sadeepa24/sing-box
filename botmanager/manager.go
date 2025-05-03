package botmanager

import (
	"errors"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/connectedbot/bottracker"
	"github.com/sagernet/sing-box/connectedbot/opts"
	botrule "github.com/sagernet/sing-box/connectedbot/rule"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/route"
)

type Manager struct {
	*bottracker.ConnManager // manage user adding removing getting status
	rules map[string]*botrule.RuleForBot // using this rule can change outbound of specific user
	router     *route.Router

}

func NewManager(inboundManager adapter.InboundManager, router *route.Router) (*Manager, error) {
	allrules := router.Rules()
	r := map[string]*botrule.RuleForBot{}
	for _, rule := range allrules {
		if rule.Type() == C.RuleTypeBot {
			r[rule.(*botrule.RuleForBot).Outbound()] = rule.(*botrule.RuleForBot)
		}
	}
	if len(r) == 0 {
		return nil, errors.New("botrule rule count cannot be zero")
	}
	connmgr :=  bottracker.NewConnManager(inboundManager)
	router.SetTracker(connmgr)
	
	return &Manager{
		ConnManager: connmgr,
		rules: r,
		router: router,
	}, nil
}

func (m *Manager) AddUser(u opts.User) (opts.UserStatus, error) {
	r, ok := m.rules[u.Outbound]
	if !ok {
		return opts.UserStatus{}, errors.New("cannot find outboun with name " + u.Outbound)
	}
	r.AddUser(u.UserStr)
	return m.ConnManager.AddUser(u)
}
func (m *Manager) AddUserReset(u opts.User) (opts.UserStatus, error) {
	m.ChangeOutbound(u)
	return m.ConnManager.AddUserReset(u)
}

func (m *Manager) RemoveUser(u opts.User) (opts.UserStatus, error) {
	m.remuserout(u)
	return m.ConnManager.RemoveUser(u)
}

func (m *Manager) ChangeOutbound(u opts.User) error {
	m.remuserout(u)
	r, ok := m.rules[u.Outbound]
	if !ok {
		return errors.New("cannot find outboun with name " + u.Outbound)
	}
	r.AddUser(u.UserStr)
	return nil
}

func (m *Manager) remuserout(u opts.User) {
	for _, r := range m.rules {
		r.RemoveUser(u.UserStr)
	}
}