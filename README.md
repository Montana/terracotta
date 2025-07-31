# Terracotta

A powerful, lightweight server tunneling program written in Go.

[![Go Version](https://img.shields.io/badge/go-1.21%2B-blue)](https://golang.org/dl/) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Overview

Terracotta is a robust and efficient TCP tunneling tool designed for developers, sysadmins, and anyone needing secure, high-performance port forwarding. With a simple CLI and three flexible operation modes, Terracotta makes it easy to expose local services, connect to remote servers, or set up a relay tunnel server.

---

## Why Terracotta?

- **Lightweight:** Minimal dependencies, fast startup, and low resource usage.
- **Flexible:** Supports direct forwarding, server, and client tunnel modes.
- **Insightful:** Real-time stats for connections, data transfer, and uptime.
- **Concurrent:** Handles multiple connections efficiently with Go's goroutines.
- **User-Friendly:** Simple command-line interface and easy configuration.

---

## Features

- **Three Operation Modes:** Direct forwarding, tunnel server, and tunnel client
- **Real-time Statistics:** Monitor connections, data transfer, and uptime
- **Secure Tunneling:** Reliable TCP tunneling with connection management
- **High Performance:** Concurrent connection handling with goroutines
- **Easy Configuration:** Simple command-line interface

---

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/terracotta.git
cd terracotta

# Build the binary
go build -o terracotta

# Or install directly (if $GOPATH/bin is in your PATH)
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

| Flag      | Description                        | Default   |
|-----------|------------------------------------|-----------|
| `-mode`     | Operation mode (`local`/`server`/`client`)| local     |
| `-local`    | Local port to listen on            | 8080      |
| `-remote`   | Remote address to connect to       | localhost |
| `-port`     | Remote port to connect to          | 80        |
| `-server`   | Server port for tunneling          | 9090      |
| `-verbose`  | Enable verbose logging             | false     |
| `-help`     | Show help information              |           |
| `-version`  | Show version information           |           |

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

## Contributing

We welcome contributions from the community! To get started:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please see our [contributing guidelines](CONTRIBUTING.md) if available.

---

## Community & Support

- Email: support@terracotta.dev
- Issues: [GitHub Issues](https://github.com/yourusername/terracotta/issues)
- Documentation: [Wiki](https://github.com/yourusername/terracotta/wiki)

If you find this project useful, please star the repository and share it with others!

