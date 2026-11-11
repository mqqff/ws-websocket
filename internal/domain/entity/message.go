package entity

import "time"

// MessageType mendefinisikan tipe dari payload pesan WebSocket
type MessageType string

const (
	MessageTypeChat   MessageType = "chat"
	MessageTypeSystem MessageType = "system"
)

// Message adalah entitas payload pesan dalam sistem live chat
type Message struct {
	Type      MessageType `json:"type"`
	Sender    string      `json:"sender"`
	Content   string      `json:"content"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewChatMessage membuat pesan obrolan reguler dari user
func NewChatMessage(sender, content string) Message {
	return Message{
		Type:      MessageTypeChat,
		Sender:    sender,
		Content:   content,
		Timestamp: time.Now(),
	}
}

// NewSystemMessage membuat pesan notifikasi sistem (misal: user join/leave)
func NewSystemMessage(content string) Message {
	return Message{
		Type:      MessageTypeSystem,
		Sender:    "System",
		Content:   content,
		Timestamp: time.Now(),
	}
}
