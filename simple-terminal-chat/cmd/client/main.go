package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		slog.Error("Error dialing:", err.Error())
		return
	}
	defer conn.Close()

	go handleConnectionRead(conn)

	message := make([]byte, 128)
	for {
		n, err := os.Stdin.Read(message)
		if err != nil {
			break
		}

		_, err = fmt.Fprint(conn, string(message[:n]))
		if err != nil {
			slog.Error(fmt.Sprintf("Error sending message: %s", err))
			break
		}
	}
}

func handleConnectionRead(conn net.Conn) {
	slog.Info("Connection established.")
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Close connection error. Error:", err.Error())
			return
		}

		slog.Info("Connection closed.")
	}(conn)

	continues := false

	for {
		message := make([]byte, 128)
		n, err := conn.Read(message)
		if err != nil {
			slog.Error(fmt.Sprintf("Error on receiving the message: %s", err.Error()))
			break
		}

		if n == 0 {
			continue
		}

		if continues {
			fmt.Printf("%s", string(message[:n]))
		} else {
			fmt.Printf("Room message: %s", string(message[:n]))
		}

		continues = message[n-1] != 10
	}
}
