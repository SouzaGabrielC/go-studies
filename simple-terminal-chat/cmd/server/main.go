package main

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log/slog"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	connections := make(map[uuid.UUID]net.Conn)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("Connection error. Error:", err.Error())
		}

		connId := uuid.New()
		connections[connId] = conn

		go handleConnectionRead(connId, conn, connections)
	}
}

func handleConnectionRead(connId uuid.UUID, conn net.Conn, connections map[uuid.UUID]net.Conn) {
	slog.Info(fmt.Sprintf("[%s] Connection established.", connId))
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error(fmt.Sprintf("[%s] Close connection error. Error: %s", connId, err.Error()))
			return
		}

		slog.Info(fmt.Sprintf("[%s] Connection closed.", connId))
	}(conn)

	message := make([]byte, 128)

	for {
		n, err := conn.Read(message)
		if err != nil {
			if err != io.EOF {
				slog.Error(fmt.Sprintf("[%s] Error on receiving the message: %s", connId, err.Error()))
			}
			break
		}

		if n == 0 {
			continue
		}

		// TODO: Change this approach to a centralized broadcast of the messages to avoid overlapping non finished messages
		broadcastMessageToConnections(connections, connId, message[:n])
	}
}

func broadcastMessageToConnections(connections map[uuid.UUID]net.Conn, senderId uuid.UUID, message []byte) {
	for connId, conn := range connections {
		if connId == senderId {
			continue
		}

		conn.Write(message)
	}
}
