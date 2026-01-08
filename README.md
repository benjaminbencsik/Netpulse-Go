# NetPulse-Go

NetPulse-Go is a simple, concurrent network monitoring tool written in Go. It reads a list of target URLs or IP addresses from a file, pings them at a regular 30-second interval, and logs the status and latency to a CSV file.

This project uses only the Go standard library.

## Installation

The easiest way to install NetPulse-Go is to use `go install`. This will download the source, compile it, and place the `Netpulse-Go` binary in your Go bin directory, making it available from your command line.

**Note:** You must have your Go environment set up correctly (i.e., your `$GOPATH/bin` or `$HOME/go/bin` directory must be in your system's `PATH`).

1.  **Install the application:**
    ```sh
    go install -v github.com/benjaminbencsik/Netpulse-Go@latest
    ```

2.  **Create your targets file:**
    `Netpulse-Go` looks for a `targets.txt` file in the directory where you run it. Create this file and add the URLs or IP addresses you want to monitor, one per line.
    
    **Example `targets.txt`:**
    ```
    https://www.google.com
    8.8.8.8
    ```

3.  **Run the application:**
    Once installed, you can run it from any directory:
    ```sh
    Netpulse-Go
    ```
    The monitor will start and create a `network_log.csv` file in the same directory. You can stop it by pressing `Ctrl+C`.

## Development (Running from Source)

If you want to modify the code or run it directly from the source, follow these steps:

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/benjaminbencsik/Netpulse-Go.git
    cd Netpulse-Go
    ```

2.  **Run the application:**
    ```sh
    go run main.go
    ```