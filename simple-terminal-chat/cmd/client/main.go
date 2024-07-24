package main

import (
	"bufio"
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

	_, err = os.Stdin.WriteTo(conn)
	if err != nil {
		slog.Error("Error writing message:", err.Error())
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
		fmt.Printf("Room message: %s\n", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
