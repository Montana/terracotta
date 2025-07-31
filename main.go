package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	version = "1.0.0"
	banner  = `
████████╗███████╗██████╗ ██████╗  █████╗  ██████╗ ██████╗ ████████╗████████╗ █████╗ 
╚══██╔══╝██╔════╝██╔══██╗██╔══██╗██╔══██╗██╔════╝██╔═══██╗╚══██╔══╝╚══██╔══╝██╔══██╗
   ██║   █████╗  ██████╔╝██████╔╝███████║██║     ██║   ██║   ██║      ██║   ███████║
   ██║   ██╔══╝  ██╔══██╗██╔══██╗██╔══██║██║     ██║   ██║   ██║      ██║   ██╔══██║
   ██║   ███████╗██║  ██║██║  ██║██║  ██║╚██████╗╚██████╔╝   ██║      ██║   ██║  ██║
   ╚═╝   ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═════╝    ╚═╝      ╚═╝   ╚═╝  ╚═╝
                                Server Tunneling v%s
`
)

type Config struct {
	Mode       string
	LocalPort  int
	RemoteAddr string
	RemotePort int
	ServerPort int
	Verbose    bool
}

type TunnelServer struct {
	config     *Config
	listener   net.Listener
	clients    map[string]net.Conn
	clientsMux sync.RWMutex
	stats      *Stats
}

type TunnelClient struct {
	config *Config
	conn   net.Conn
	stats  *Stats
}

type Stats struct {
	ConnectTime   time.Time
	BytesSent     int64
	BytesReceived int64
	ActiveConns   int64
	TotalConns    int64
	mu            sync.RWMutex
}

func (s *Stats) AddConnection() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ActiveConns++
	s.TotalConns++
}

func (s *Stats) RemoveConnection() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ActiveConns--
}

func (s *Stats) AddBytes(sent, received int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.BytesSent += sent
	s.BytesReceived += received
}

func (s *Stats) GetStats() (int64, int64, int64, int64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ActiveConns, s.TotalConns, s.BytesSent, s.BytesReceived
}

func main() {
	config := parseFlags()
	fmt.Printf(banner, version)
	fmt.Printf("Starting Terracotta in %s mode...\n\n", config.Mode)
	switch config.Mode {
	case "server":
		runServer(config)
	case "client":
		runClient(config)
	case "local":
		runLocalTunnel(config)
	default:
		log.Fatal("Invalid mode. Use 'server', 'client', or 'local'")
	}
}

func parseFlags() *Config {
	config := &Config{}
	flag.StringVar(&config.Mode, "mode", "local", "Mode: server, client, or local")
	flag.IntVar(&config.LocalPort, "local", 8080, "Local port to listen on")
	flag.StringVar(&config.RemoteAddr, "remote", "localhost", "Remote address to connect to")
	flag.IntVar(&config.RemotePort, "port", 80, "Remote port to connect to")
	flag.IntVar(&config.ServerPort, "server", 9090, "Server port for tunneling")
	flag.BoolVar(&config.Verbose, "verbose", false, "Verbose logging")
	showHelp := flag.Bool("help", false, "Show help")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()
	if *showVersion {
		fmt.Printf("Terracotta v%s\n", version)
		os.Exit(0)
	}
	if *showHelp {
		printHelp()
		os.Exit(0)
	}
	return config
}

func printHelp() {
	fmt.Printf(banner, version)
	fmt.Println(`
Usage: terracotta [options]

Modes:
  local   - Direct port forwarding (default)
  server  - Run as tunnel server
  client  - Connect to tunnel server

Options:
  -mode string      Mode: server, client, or local (default "local")
  -local int        Local port to listen on (default 8080)
  -remote string    Remote address to connect to (default "localhost")
  -port int         Remote port to connect to (default 80)
  -server int       Server port for tunneling (default 9090)
  -verbose          Enable verbose logging
  -help             Show this help
  -version          Show version

Examples:
  # Direct port forwarding
  terracotta -local 8080 -remote example.com -port 80
  
  # Run tunnel server
  terracotta -mode server -server 9090
  
  # Run tunnel client
  terracotta -mode client -local 8080 -remote tunnelserver.com -server 9090
  `)
}

func runLocalTunnel(config *Config) {
	stats := &Stats{ConnectTime: time.Now()}
	localAddr := fmt.Sprintf(":%d", config.LocalPort)
	remoteAddr := fmt.Sprintf("%s:%d", config.RemoteAddr, config.RemotePort)
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", localAddr, err)
	}
	defer listener.Close()
	fmt.Printf("Terracotta tunnel active\n")
	fmt.Printf("Local: %s -> Remote: %s\n", localAddr, remoteAddr)
	fmt.Printf("Forwarding traffic...\n\n")
	go printStats(stats, config.Verbose)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\nShutting down Terracotta...\n")
		listener.Close()
		os.Exit(0)
	}()
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			if config.Verbose {
				log.Printf("Accept error: %v", err)
			}
			continue
		}
		stats.AddConnection()
		go handleLocalConnection(clientConn, remoteAddr, stats, config.Verbose)
	}
}

func handleLocalConnection(clientConn net.Conn, remoteAddr string, stats *Stats, verbose bool) {
	defer func() {
		clientConn.Close()
		stats.RemoveConnection()
	}()
	if verbose {
		log.Printf("New connection from %s", clientConn.RemoteAddr())
	}
	remoteConn, err := net.DialTimeout("tcp", remoteAddr, 10*time.Second)
	if err != nil {
		if verbose {
			log.Printf("Failed to connect to %s: %v", remoteAddr, err)
		}
		return
	}
	defer remoteConn.Close()
	done := make(chan struct{}, 2)
	go func() {
		defer func() { done <- struct{}{} }()
		sent, _ := io.Copy(remoteConn, clientConn)
		stats.AddBytes(sent, 0)
		if verbose {
			log.Printf("Client->Remote finished: %d bytes", sent)
		}
	}()
	go func() {
		defer func() { done <- struct{}{} }()
		received, _ := io.Copy(clientConn, remoteConn)
		stats.AddBytes(0, received)
		if verbose {
			log.Printf("Remote->Client finished: %d bytes", received)
		}
	}()
	<-done
}

func runServer(config *Config) {
	server := &TunnelServer{
		config:  config,
		clients: make(map[string]net.Conn),
		stats:   &Stats{ConnectTime: time.Now()},
	}
	serverAddr := fmt.Sprintf(":%d", config.ServerPort)
	listener, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Failed to start server on %s: %v", serverAddr, err)
	}
	defer listener.Close()
	server.listener = listener
	fmt.Printf("Terracotta server started\n")
	fmt.Printf("Listening on: %s\n", serverAddr)
	fmt.Printf("Waiting for clients...\n\n")
	go printStats(server.stats, config.Verbose)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\nShutting down server...\n")
		server.shutdown()
		os.Exit(0)
	}()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if config.Verbose {
				log.Printf("Accept error: %v", err)
			}
			continue
		}
		server.stats.AddConnection()
		go server.handleClient(conn)
	}
}

func (s *TunnelServer) handleClient(conn net.Conn) {
	defer func() {
		conn.Close()
		s.stats.RemoveConnection()
	}()
	clientAddr := conn.RemoteAddr().String()
	if s.config.Verbose {
		log.Printf("New client connected: %s", clientAddr)
	}
	s.clientsMux.Lock()
	s.clients[clientAddr] = conn
	s.clientsMux.Unlock()
	defer func() {
		s.clientsMux.Lock()
		delete(s.clients, clientAddr)
		s.clientsMux.Unlock()
	}()
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if s.config.Verbose {
				log.Printf("Client %s disconnected: %v", clientAddr, err)
			}
			break
		}
		s.stats.AddBytes(0, int64(n))
		conn.Write(buffer[:n])
		s.stats.AddBytes(int64(n), 0)
	}
}

func (s *TunnelServer) shutdown() {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()
	for addr, conn := range s.clients {
		if s.config.Verbose {
			log.Printf("Closing connection to %s", addr)
		}
		conn.Close()
	}
	if s.listener != nil {
		s.listener.Close()
	}
}

func runClient(config *Config) {
	client := &TunnelClient{
		config: config,
		stats:  &Stats{ConnectTime: time.Now()},
	}
	serverAddr := fmt.Sprintf("%s:%d", config.RemoteAddr, config.ServerPort)
	localAddr := fmt.Sprintf(":%d", config.LocalPort)
	fmt.Printf("Terracotta client starting\n")
	fmt.Printf("Server: %s\n", serverAddr)
	fmt.Printf("Local: %s\n", localAddr)
	fmt.Printf("Establishing tunnel...\n\n")
	conn, err := net.DialTimeout("tcp", serverAddr, 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to server %s: %v", serverAddr, err)
	}
	defer conn.Close()
	client.conn = conn
	fmt.Printf("Connected to tunnel server\n")
	go printStats(client.stats, config.Verbose)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\nDisconnecting from server...\n")
		conn.Close()
		os.Exit(0)
	}()
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", localAddr, err)
	}
	defer listener.Close()
	fmt.Printf("Tunnel established! Local port %d is now forwarded\n", config.LocalPort)
	for {
		localConn, err := listener.Accept()
		if err != nil {
			if config.Verbose {
				log.Printf("Accept error: %v", err)
			}
			continue
		}
		client.stats.AddConnection()
		go client.handleLocalConnection(localConn)
	}
}

func (c *TunnelClient) handleLocalConnection(localConn net.Conn) {
	defer func() {
		localConn.Close()
		c.stats.RemoveConnection()
	}()
	if c.config.Verbose {
		log.Printf("New local connection from %s", localConn.RemoteAddr())
	}
	done := make(chan struct{}, 2)
	go func() {
		defer func() { done <- struct{}{} }()
		sent, _ := io.Copy(c.conn, localConn)
		c.stats.AddBytes(sent, 0)
	}()
	go func() {
		defer func() { done <- struct{}{} }()
		received, _ := io.Copy(localConn, c.conn)
		c.stats.AddBytes(0, received)
	}()
	<-done
}

func printStats(stats *Stats, verbose bool) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		active, total, sent, received := stats.GetStats()
		uptime := time.Since(stats.ConnectTime).Round(time.Second)
		if verbose || active > 0 {
			fmt.Printf("Stats - Active: %d, Total: %d, Sent: %s, Received: %s, Uptime: %s\n",
				active, total, formatBytes(sent), formatBytes(received), uptime)
		}
	}
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
