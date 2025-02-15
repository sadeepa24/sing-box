package box

import (
	"errors"
	"fmt"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/connectedbot"
)

func (b *Box) AddUser(BotUser connectedbot.BotUser, inboundtag string) (connectedbot.StatusOutput, error) {
	inb, found := b.inbound.Get(inboundtag)
	if !found {
		return connectedbot.StatusOutput{}, fmt.Errorf("inbound not found")
	}
	botFace, ok := inb.(adapter.BotFace)
	if !ok {
		return connectedbot.StatusOutput{}, fmt.Errorf("inbound does not support")
	}
	return botFace.AddUser(BotUser)
}

func (b *Box) RemoveUser(BotUser connectedbot.BotUser, inboundtag string) (connectedbot.StatusOutput, error) {
	inb, found := b.inbound.Get(inboundtag)
	if !found {
		return connectedbot.StatusOutput{}, fmt.Errorf("inbound not found")
	}
	botFace, ok := inb.(adapter.BotFace)
	if !ok {
		return connectedbot.StatusOutput{}, fmt.Errorf("inbound does not support")
	}
	return botFace.RemoveUser(BotUser)
	

}
func (b *Box) GetastatusUser(BotUser connectedbot.BotUser, inboundtag string) (connectedbot.StatusOutput, error) {
	inb, found := b.inbound.Get(inboundtag)
	if !found {
		return connectedbot.StatusOutput{}, fmt.Errorf("inbound not found")
	}
	botFace, ok := inb.(adapter.BotFace)
	if !ok {
		return connectedbot.StatusOutput{}, fmt.Errorf("inbound does not support")
	}
	return botFace.GetastatusUser(BotUser)

}

func (b *Box) AddUserReset(BotUser connectedbot.BotUser, inboundtag string) (connectedbot.StatusOutput, error) {
	inb, found := b.inbound.Get(inboundtag)
	if !found {
		return connectedbot.StatusOutput{}, fmt.Errorf("inbound not found")
	}
	botFace, ok := inb.(adapter.BotFace)
	if !ok {
		return connectedbot.StatusOutput{}, fmt.Errorf("inbound does not support")
	}
	return botFace.AddUserReset(BotUser)
	//return connectedbot.StatusOutput{}, vless.ErrInboundNotFound

}

func (b *Box) CloseAll(BotUser connectedbot.BotUser, inboundtag string)  error {
	
	inb, found := b.inbound.Get(inboundtag)
	if !found {
		return fmt.Errorf("inbound not found")
	}
	botFace, ok := inb.(adapter.BotFace)
	if !ok {
		return fmt.Errorf("inbound does not support")
	}
	return botFace.CloseAll(BotUser)

}



func (b *Box) Addoutbounduser(user, outbound string) {
	// adder ,ok := b.router.(adapter.OutboundAdder)
	b.router.AddOutboundUser(user, outbound)
	
	// if ok {
	// 	adder.AddOutboundUser(user, outbound)
	// }
}

func (b *Box) RemoveAllRule(user string) {
	// adder ,ok := b.router.(adapter.OutboundAdder)
	// if ok {
	// 	adder.RemoveAllRule(user)
	// }
	b.router.RemoveAllRule(user)
}

func (b *Box) GetOutBound(tag string) (adapter.Outbound, error) {
	
	ot, found := b.outbound.Outbound(tag)
	if !found {
		return nil, errors.New("outbound not found")
	}
	//b.endpoint.

	return ot, nil

}