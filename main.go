package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

func main() {
	// Define standard flags with defaults suited for containers
	start := flag.Int("start", 8000, "Start port range")
	end := flag.Int("end", 8080, "End port range")
	flag.Parse()

	// Loop through your custom range
	for port := *start; port <= *end; port++ {
		address := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))

		// Try to claim the port. If successful, it is completely free!
		listener, err := net.Listen("tcp", address)
		if err == nil {
			listener.Close() // Release it instantly so your container can take it
			fmt.Printf("%d\n", port) // Print just the port number for easy script usage
			os.Exit(0)
		}
	}

	fmt.Println("Error: No free ports found in range.")
	os.Exit(1)
}
