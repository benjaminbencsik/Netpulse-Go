# NetPulse-Go

NetPulse-Go is a simple, concurrent network monitoring tool written in Go. It reads a list of target URLs or IP addresses from a file, pings them at a regular 30-second interval, and logs the status and latency to a CSV file.

This project uses only the Go standard library.

## Features

- **Concurrent Monitoring**: Uses Goroutines and WaitGroups to check multiple targets simultaneously without blocking.
- **Configurable Targets**: Easily add or remove targets by editing the `targets.txt` file.
- **Periodic Checks**: Pings all targets every 30 seconds.
- **CSV Logging**: Records results with a timestamp, target, status, and latency in `network_log.csv`.
- **Standard Library Only**: Built using only the Go standard library, with no external dependencies.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or higher is recommended).
- A correctly configured Go environment (i.e., the `go` command is in your system's PATH).

## Getting Started

1.  **Place the project files** in a directory of your choice.

2.  **Add your targets:**
    Open the `targets.txt` file and add the URLs (e.g., `https://www.google.com`) or IP addresses/hostnames (e.g., `8.8.8.8`) you wish to monitor. Add one entry per line.

3.  **Run the application:**
    Open your terminal, navigate to the project directory, and run the following commands:

    ```sh
    # Initialize the Go module (only needs to be done once)
    go mod init netpulse-go

    # Run the network monitor
    go run main.go
    ```

    The application will start, and you will see output in your console as it begins to ping the targets. To stop the monitor, press `Ctrl+C`.

## Log Output

The monitoring results will be saved to the `network_log.csv` file.

- **Timestamp**: The UTC timestamp of the check in RFC3339 format.
- **Target**: The URL or IP address that was checked.
- **Status**: The result of the check.
    - For HTTP/S targets, this will be the HTTP status (e.g., `200 OK`).
    - For IP/hostname targets, this will be `UP` or `DOWN`.
- **Latency**: The time it took to get a response from the target.

### Example Log (`network_log.csv`)
```csv
Timestamp,Target,Status,Latency
2026-01-08T18:30:00Z,https://www.google.com,200 OK,25.123ms
2026-01-08T18:30:00Z,8.8.8.8,UP,12.456ms
2026-01-08T18:30:00Z,https://thisurldoesnotexist12345.com,Error: Get "https://thisurldoesnotexist12345.com": dial tcp: no such host,5.2s
```
