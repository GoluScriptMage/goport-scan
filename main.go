package main

import (
	"fmt"
	"net"
	"time"
)

type ScanResult struct {
	PORT   int
	IsOpen bool
	Err    error
}

func main() {
	println("Starting the process...")
	result := scanPort(22)
	fmt.Printf("Port: %d | IsOpen: %t | Error: %v\n", result.PORT, result.IsOpen, result.Err)
}

func scanPort(port int) ScanResult {

	address := fmt.Sprintf("scanme.nmap.org:%d", port)

	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		println("Error: problem in connecting to port", err)
		return ScanResult{PORT: port, IsOpen: false, Err: err}
	}
	defer conn.Close()

	return ScanResult{PORT: port, IsOpen: true, Err: nil}
}
