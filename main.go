package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// PingResult holds the data from a single ping attempt.
type PingResult struct {
	Timestamp string
	Target    string
	Status    string
	Latency   time.Duration
}

func main() {
	// Start the monitoring process
	monitor("targets.txt", "network_log.csv")
}

// monitor orchestrates the network monitoring.
func monitor(targetsFile, logFile string) {
	// Open and prepare the log file for writing.
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file '%s': %v", logFile, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write a header to the CSV file if it's new/empty.
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		if err := writer.Write([]string{"Timestamp", "Target", "Status", "Latency"}); err != nil {
			log.Printf("Failed to write header to log file: %v", err)
		}
		writer.Flush()
	}

	resultsChan := make(chan PingResult)
	var wg sync.WaitGroup

	// Start a single goroutine to handle writing all results to the log file.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range resultsChan {
			record := []string{
				result.Timestamp,
				result.Target,
				result.Status,
				result.Latency.String(),
			}
			if err := writer.Write(record); err != nil {
				log.Printf("Failed to write record to log file: %v", err)
			}
			writer.Flush() // Ensure data is written promptly.
		}
	}()

	// Perform an initial round of pings before starting the ticker.
	performPings(targetsFile, resultsChan)

	// Set up a ticker to ping all targets every 30 seconds.
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		performPings(targetsFile, resultsChan)
	}

	close(resultsChan)
	wg.Wait() // Wait for the logger goroutine to finish.
}

// performPings reads targets and initiates a concurrent ping for each.
func performPings(targetsFile string, resultsChan chan<- PingResult) {
	log.Println("Pinging targets...")
	targets, err := readTargets(targetsFile)
	if err != nil {
		log.Printf("Error reading targets file: %v", err)
		return
	}

	var pingWg sync.WaitGroup
	for _, target := range targets {
		pingWg.Add(1)
		go pingTarget(target, &pingWg, resultsChan)
	}

	pingWg.Wait()
	log.Println("Finished pinging targets.")
}

// readTargets reads the target URLs/IPs from the specified file.
func readTargets(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var targets []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			targets = append(targets, line)
		}
	}
	return targets, scanner.Err()
}

// pingTarget checks a single target and sends the result to a channel.
func pingTarget(target string, wg *sync.WaitGroup, resultsChan chan<- PingResult) {
	defer wg.Done()

	var status string
	var latency time.Duration
	startTime := time.Now()

	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		// Handle HTTP/S URLs
		resp, err := http.Get(target)
		latency = time.Since(startTime)
		if err != nil {
			status = fmt.Sprintf("Error: %s", err)
		} else {
			status = resp.Status
			resp.Body.Close()
		}
	} else {
		// Handle IP addresses or hostnames with a default TCP dial to port 80.
		conn, err := net.DialTimeout("tcp", target+":80", 5*time.Second)
		latency = time.Since(startTime)
		if err != nil {
			status = "DOWN"
		} else {
			status = "UP"
			conn.Close()
		}
	}

	resultsChan <- PingResult{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Target:    target,
		Status:    status,
		Latency:   latency,
	}
}
