package adapter

import "github.com/sagernet/sing-box/connectedbot"

type BotFace interface {
	GetastatusUser(connectedbot.BotUser) (connectedbot.StatusOutput, error)
	AddUserReset(connectedbot.BotUser) (connectedbot.StatusOutput, error)
	RemoveUser(connectedbot.BotUser) (connectedbot.StatusOutput, error)
	AddUser(connectedbot.BotUser) (connectedbot.StatusOutput, error)
	CloseAll(connectedbot.BotUser) error
}