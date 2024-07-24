package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"simple-terminal-chat/internals/message"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		slog.Error("Error dialing:", err.Error())
		return
	}
	defer conn.Close()

	fmt.Print("Enter your username: ")
	clientName, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		slog.Error("Error getting username:", err.Error())
		return
	}

	go handleConnectionRead(conn)

	scanner := bufio.NewScanner(os.Stdin)
	jsonEncoderConn := json.NewEncoder(conn)

	for scanner.Scan() {
		text := scanner.Text()

		err := jsonEncoderConn.Encode(message.Content{
			Username: clientName[:len(clientName)-2],
			Message:  text,
		})
		if err != nil {
			slog.Error("Error writing to server:", err.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
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

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		var msg message.Content

		err := json.Unmarshal(scanner.Bytes(), &msg)
		if err != nil {
			slog.Error("Error unmarshalling message:", err.Error())
			continue
		}

		fmt.Printf("Room message from [%s]: %s\n", msg.Username, msg.Message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
