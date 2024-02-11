# Command Execution API

A Go-based web api for executing shell commands via HTTP requests. It returns the command's output (stdout and stderr) and an exit code.

## Getting Started

### Prerequisites

- Go (1.21 or newer)

### Installation

- Clone the repository:

```
  git clone https://github.com/yuktea/golang-d.git
  cd golang-d
```

- (Optional) Configure the server port by creating a .env file:
- Start the server with:

  ```
  go run main.go
  ```
- Usage:
  ```
  curl -X POST 'http://localhost:8080/api/cmd' \
  -H 'Content-Type: application/json' \
  -d '{"command":"echo Hello World"}'
  ```
  ---

