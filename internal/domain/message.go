package domain

type MessageType string

const (
	TypeChat     MessageType = "chat"
	TypePresence MessageType = "presence"
)

type Message struct {
	Type    MessageType `json:"type"`
	Content any         `json:"content"`
	UserName string `json:"user_name"`
}

type ChatConten struct {
	Chat string `json:"chat"`
}

type PresenceContent struct {
	OnlineUsers []string `json:"online_users"`
}
