package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

type ScanResult struct {
	PORT    int
	IsOpen  bool
	Err     error
	Service string
}

func main() {

	// Step 1. Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	ports := make(chan int, 100)          // Channel to send ports
	results := make(chan ScanResult, 100) // Channel to receive results

	// Flags
	host := flag.String("host", "scanme.nmap.org", "Host to scan")
	start := flag.Int("start", 1, "Start port")
	end := flag.Int("end", 1024, "End port")
	workers := flag.Int("workers", 100, "Number of workers")
	flag.Parse() // Use this after defining the flags

	// Step 2. Spawn goroutines to handle results
	startTime := time.Now()
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go worker(ports, results, &wg, *host)
	}

	// Step 3. Send ports to port channel in a separate goroutine
	go func() {
		for i := *start; i <= *end; i++ {
			ports <- i
		}
		close(ports) // Close the ports channel after sending all ports
	}()

	// Step 4. Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results) // Close the results channel after all workers are done
	}()

	// Step 5. Collect results from the results channel
	for res := range results {
		if res.IsOpen {
			fmt.Printf("Port %d is open\n", res.PORT)
			fmt.Printf("Service %s%d is open\n", res.Service, res.PORT)
		}
	}
	fmt.Printf("Scan completed in %v\n", time.Since(startTime))
}

// Scans a single port and returns the result
func scanPort(port int, host string) ScanResult {

	address := fmt.Sprintf("%s:%d", host, port)
	// address := fmt.Sprintf("127.0.0.1:%d", port)

	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return ScanResult{PORT: port, IsOpen: false, Err: err, Service: ""}
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(2 * time.Second)) // Set deadline of after 2 sec
	buffer := make([]byte, 256)
	n, readErr := conn.Read(buffer)
	var banner string

	if readErr != nil {
		banner = "Unknown/Silent"
	} else {
		banner = string(buffer[:n]) // n means from 0:7 index numbers from 256 slots arr
	}

	return ScanResult{PORT: port, IsOpen: true, Err: nil, Service: banner}
}

// Woker function that scans ports received from the ports channel and sends results to the results channel
func worker(ports chan int, results chan ScanResult, wg *sync.WaitGroup, host string) {

	// Scan the ports received from the ports channel
	for p := range ports {
		scanResult := scanPort(p, host)
		// Send the scan results to results channel
		results <- scanResult
	}
	wg.Done() // Mark the worker as done when all ports are scanned
}
