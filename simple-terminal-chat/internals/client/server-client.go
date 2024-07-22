package client

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net"
	"simple-terminal-chat/internals/message"
)

type ConnectionStatus uint8

const (
	ConnectionOpen ConnectionStatus = iota
	ConnectionClosed
)

type ServerClient struct {
	Id               uuid.UUID
	connection       net.Conn
	connectionStatus ConnectionStatus
}

func NewServerClient(conn net.Conn) *ServerClient {
	return &ServerClient{
		Id:               uuid.New(),
		connection:       conn,
		connectionStatus: ConnectionOpen,
	}
}

func (c *ServerClient) Close() {
	err := c.connection.Close()
	if err != nil {
		slog.Error(fmt.Sprintf("[%s] Close connection error. Error: %s", c.Id, err.Error()))
	}

	slog.Info(fmt.Sprintf("[%s] Connection closed.", c.Id))
	c.connectionStatus = ConnectionClosed
}

func (c *ServerClient) ReadMessage() *message.Message {
	messageContentBytes := make([]byte, 128)

	pr, pw := io.Pipe()

	receivedMessage := &message.Message{
		MessageReader: pr,
		SenderId:      c.Id,
	}

	// Loop to receive chunks of the same message and pipe it
	go func() {
		// Close pipe when finish reading loop
		defer pw.Close()

		for {
			n, err := c.connection.Read(messageContentBytes)
			if err != nil {
				if err == io.EOF {
					c.Close()
				} else {
					slog.Error(fmt.Sprintf("[%s] Error on receiving the messageContentBytes: %s", c.Id, err.Error()))
				}

				break
			}

			if n == 0 {
				break
			}

			pw.Write(messageContentBytes[:n])
		}
	}()

	return receivedMessage
}

func (c *ServerClient) Write(message []byte) {
	c.connection.Write(message)
}
