package opts

import (
	"fmt"

	"github.com/sagernet/sing-box/option"
	M "github.com/sagernet/sing/common/metadata"
)

type User struct {
	MaxLogin  int16
	Bandwidth int64 //byte
	UserStr   string
	Uid       int
	Outbound  string
	Proto ComProto
	InboundList []string
}

type UserStatus struct {
	Download int64
	Upload   int64
	Ip       map[string]int16
	Disabled bool

	UserID	int
}

type CallBackResult struct {
	Status UserStatus
	Source      M.Socksaddr
	Destination M.Socksaddr
	Protocol     string
	Domain       string
	Outbound    string
	Inbound     string
	User        string
	Network string
}

func (u *CallBackResult) String() string {
	return fmt.Sprintf(
		"Connection Info:\n"+
			"  From      : %s\n"+
			"  To        : %s\n"+
			"  Protocol  : %s\n"+
			"  Domain    : %s\n"+
			"  Inbound   : %s\n"+
			"  Outbound  : %s\n"+
			"  User      : %s\n"+
			"  Network   : %s\n"+
			"  Download  : %d bytes\n"+
			"  Upload    : %d bytes",
		u.Source.String(),
		u.Destination.String(),
		u.Protocol,
		u.Domain,
		u.Inbound,
		u.Outbound,
		u.User,
		u.Network,
		u.Status.Download,
		u.Status.Upload,
	)
}

type ComProto interface {
	Vless() (option.VLESSUser, bool)
	Vmess() (option.VMessUser, bool)
	Trojan() (option.TrojanUser, bool)
	UserStr() string
	Uid() int
	Password() string
	UUID() string
} 