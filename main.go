package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type ScanResult struct {
	PORT   int
	IsOpen bool
	Err    error
}

func main() {

	// Step 1. Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	ports := make(chan int, 100)          // Channel to send ports
	results := make(chan ScanResult, 100) // Channel to receive results

	// Step 2. Spawn goroutines to handle results
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go worker(ports, results, &wg) 
	}

	// Step 3. Send ports to port channel in a separate goroutine
	go func() {
		for i := 1; i <= 1024; i++ {
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
		}
	}
}

// Scans a single port and returns the result
func scanPort(port int) ScanResult {

	address := fmt.Sprintf("scanme.nmap.org:%d", port)

	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return ScanResult{PORT: port, IsOpen: false, Err: err}
	}
	defer conn.Close()

	return ScanResult{PORT: port, IsOpen: true, Err: nil}
}

// Woker function that scans ports received from the ports channel and sends results to the results channel
func worker(ports chan int, results chan ScanResult, wg *sync.WaitGroup) {

	// Scan the ports received from the ports channel
	for p := range ports {
		scanResult := scanPort(p)
		// Send the scan results to results channel
		results <- scanResult
	}
	wg.Done() // Mark the worker as done when all ports are scanned
}
