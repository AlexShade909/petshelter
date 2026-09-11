package network

import (
	"fmt"
	"net"
)

func Connect() (net.Conn, error) {
	fmt.Println("Start server")
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error connection", err)
		return nil, err
	}
	defer listener.Close()
	fmt.Println("Wait connection client...")
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error connection client", err)
		return nil, err
	}
	fmt.Println("Connection success. Listen more...")

	return conn, nil
}
