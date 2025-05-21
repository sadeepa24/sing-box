package opts

import (
	"github.com/sagernet/sing-box/option"
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

type ComProto interface {
	Vless() (option.VLESSUser, bool)
	Vmess() (option.VMessUser, bool)
	Trojan() (option.TrojanUser, bool)
	UserStr() string
	Uid() int
	Password() string
	UUID() string
} 