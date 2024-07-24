package message

import (
	"github.com/google/uuid"
)

type Content struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Message struct {
	SenderId       uuid.UUID
	MessageContent Content
}
