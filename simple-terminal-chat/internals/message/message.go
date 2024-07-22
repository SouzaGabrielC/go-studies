package message

import (
	"github.com/google/uuid"
	"io"
)

type Message struct {
	MessageReader io.Reader
	SenderId      uuid.UUID
}
