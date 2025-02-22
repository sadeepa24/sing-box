package route

import (
	"context"
	"sync"

	"github.com/sagernet/sing-box/adapter"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing-box/route/rule"
)

func (r *Router) AddOutboundUser(user, outboundtag string) {
	r.RemoveAllRule(user)
	
	for _, rule := range r.rules {
		if ruler, ok := rule.(Botrule);ok {
			if ruler.Outbound() == outboundtag {
				ruler.AddUser(user)
				break
			}
		}
	}

}


func (r *Router) RemoveAllRule(user string) {
	for _, rule := range r.rules {
		if ruler, ok := rule.(Botrule); ok {
			ruler.RemoveUser(user)
		}
		
	}

}

// func (r *Router) SkipSniff(user string ) bool {
// 	_, ok := r.directusermap.Load(user)
// 	return ok
// }



type Botrule interface {
	AddUser(user string)
	RemoveUser(user string)
	Outbound() string
}


type RuleForBot struct {
	//users *sync.Map
	usermap map[string]bool
	mu *sync.RWMutex
	action adapter.RuleAction
	outbound string
}


func NewBotRule(ctx context.Context, logger log.ContextLogger, options option.Rule, checkOutbound bool)(adapter.Rule, error) {
	action, err := rule.NewRuleAction(ctx, logger, options.DefaultOptions.RuleAction)
	if err != nil {
		return nil, err
	}
	return &RuleForBot{
		//users: &sync.Map{},
		mu: &sync.RWMutex{},
		usermap: map[string]bool{},
		action: action,
		outbound: options.DefaultOptions.RuleAction.RouteOptions.Outbound,
		
		

	}, nil
}


func (r *RuleForBot) Start() error {return nil}
func (r *RuleForBot) Close() error {return nil}
func (r *RuleForBot) Action() adapter.RuleAction { return r.action }


func (r *RuleForBot) Match(metadata *adapter.InboundContext) bool {
	r.mu.RLock()
	_, ok := r.usermap[metadata.User]
	r.mu.RUnlock()
	return ok
}

func (r *RuleForBot) AddUser(user string) {

	r.mu.Lock()
	r.usermap[user] = true
	r.mu.Unlock()
}

func (r *RuleForBot) RemoveUser(user string) {
	r.mu.Lock()
	delete(r.usermap, user)
	r.mu.Unlock()
}


func (r *RuleForBot) String() string { 
	return "rule only for connected bot " + r.outbound
}


func (r *RuleForBot) Type() string { return C.RuleTypeDefault }
func (r *RuleForBot) UpdateGeosite() error {return nil }
func (r *RuleForBot) Outbound() string {return r.outbound }