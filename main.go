package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "4040"
	CONN_TYPE = "tcp"
	VERSION = "Prerun v0.1a"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	// Listen for an incoming connection.
	conn, err := l.Accept()
	if err != nil {
		fmt.Printf("Error accepting: ", err.Error())
		os.Exit(1)
	}
	// Handle connections in a new goroutine.
	handleRequest(conn)

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	fmt.Println("Handling the request!")
	// Make a buffer to hold xfer data.
	buf_receive := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	for {
		msgLen, err := receivePacket(conn, buf_receive)
		checkReadError(err)
		version := string(buf_receive[:msgLen])
		if version == VERSION {
			//we start the interface.
			cmd := displayMenu()
			fmt.Printf("Command: %q\n", cmd)
		} else {
			fmt.Errorf("Wrong version reveived! Expected: %q , Received: %q", VERSION, version)
			break
		}
	}
	conn.Close()
}

func displayMenu() (string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Menu:")
	fmt.Println("(1) pwd")
	fmt.Println("(2) ls")
	fmt.Println("(q) Quit")
	input, err := reader.ReadString('\n')
	checkReadError(err)
	if len(input) > 1 {
		input = input[:len(input) - 1]
	}
	fmt.Printf("input :%q", input)
	checkReadError(err)
	var cmd string
	switch input {
		case "1":
			cmd = "pwd"
		case "2":
			cmd = "ls"
		default:
		case "q":
			cmd = "quit"
	}
	return cmd
}

func receivePacket(conn net.Conn, buf_receive []byte) (int, error) {
	reqLen, err := conn.Read(buf_receive)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return reqLen, err
	} else if reqLen != 0 {
		fmt.Printf("Received: %q", buf_receive[:reqLen])

		// Send a response back to person contacting us.
		// Include package number!!
		conn.Write([]byte("ACK"))
		// Close the connection when you're done with it.
	}
	return reqLen, nil
}

func checkReadError(err error) {
	if err != nil {
		fmt.Errorf("Error: %q", err.Error())
		os.Exit(1)
	}
}