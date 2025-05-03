package botrule

import (
	"context"
	"sync"

	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/route/rule"
)


type Botrule interface {
	AddUser(user string)
	RemoveUser(user string)
	Outbound() string
}


type RuleForBot struct {
	//users *sync.Map
	usermap sync.Map
	action adapter.RuleAction
	outbound string
}


func NewBotRule(ctx context.Context, logger log.ContextLogger, options option.Rule, checkOutbound bool)(adapter.Rule, error) {
	action, err := rule.NewRuleAction(ctx, logger, options.DefaultOptions.RuleAction)
	if err != nil {
		return nil, err
	}
	return &RuleForBot{
		action: action,
		outbound: options.DefaultOptions.RuleAction.RouteOptions.Outbound,
	}, nil
}


func (r *RuleForBot) Start() error {return nil}
func (r *RuleForBot) Close() error {return nil}
func (r *RuleForBot) Action() adapter.RuleAction { return r.action }


func (r *RuleForBot) Match(metadata *adapter.InboundContext) bool {
	_, ok := r.usermap.Load(metadata.User)
	return ok
}

func (r *RuleForBot) AddUser(user string) {
	r.usermap.Store(user, true)
}

func (r *RuleForBot) RemoveUser(user string) {
	r.usermap.Delete(user)
}


func (r *RuleForBot) String() string { 
	return "rule only for connected bot " + r.outbound
}


func (r *RuleForBot) Type() string { return C.RuleTypeBot }
func (r *RuleForBot) UpdateGeosite() error {return nil }
func (r *RuleForBot) Outbound() string {return r.outbound }