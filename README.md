# terracotta

terracotta is a lightweight server tunneling program written in Go that i made for myself now i'm publishing it

<img width="442" height="119" alt="Screenshot 2025-07-31 at 2 13 49 PM" src="https://github.com/user-attachments/assets/ad985dfa-685e-4483-aed1-faa1e1c2b436" />

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://golang.org/dl/) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

### Installation

```bash
git clone https://github.com/montana/terracotta.git
cd terracotta
go build -o terracotta
go install
```

---

## Usage

### Direct Port Forwarding
Forward local port 8080 to example.com:80:
```bash
./terracotta -local 8080 -remote example.com -port 80
```

### Tunnel Server
Run a tunnel server on port 9090:
```bash
./terracotta -mode server -server 9090
```

### Tunnel Client
Connect to tunnel server and forward local port:
```bash
./terracotta -mode client -local 8080 -remote tunnelserver.com -server 9090
```

---

## Command Line Options

| Flag        | Description                                 | Default   |
|-------------|---------------------------------------------|-----------|
| `-mode`     | Operation mode (`local`/`server`/`client`)  | local     |
| `-local`    | Local port to listen on                     | 8080      |
| `-remote`   | Remote address to connect to                | localhost |
| `-port`     | Remote port to connect to                   | 80        |
| `-server`   | Server port for tunneling                   | 9090      |
| `-verbose`  | Enable verbose logging                      | false     |
| `-help`     | Show help information                       |           |
| `-version`  | Show version information                    |           |

---

## Examples

### Web Development
Forward your local development server to an external service:
```bash
./terracotta -local 3000 -remote api.example.com -port 443
```

### Database Access
Tunnel to a remote database through a bastion host:
```bash
./terracotta -local 5432 -remote db.internal.com -port 5432
```

### Expose Local Service
Make a local service available through a tunnel server:
```bash
# On the server
./terracotta -mode server -server 9090

# On the client
./terracotta -mode client -local 8080 -remote your-server.com -server 9090
```

---

## Building from Source

### Prerequisites
- Go 1.21 or later

### Build
```bash
go build -o terracotta
```

### Cross-compilation
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o terracotta-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o terracotta.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o terracotta-macos
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Author

Michael Mendy © 2025
