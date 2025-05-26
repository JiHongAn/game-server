package socket

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type MatchServer struct {
	listener   net.Listener
	clients    map[string]net.Conn
	clientsMux sync.RWMutex
}

func NewMatchServer() *MatchServer {
	return &MatchServer{
		clients: make(map[string]net.Conn),
	}
}

func (s *MatchServer) Start(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start match server: %v", err)
	}
	s.listener = listener

	log.Printf("Match Server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *MatchServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *MatchServer) handleConnection(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	log.Printf("New client connected: %s", clientAddr)

	s.addClient(clientAddr, conn)
	defer s.removeClient(clientAddr)
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from client %s: %v", clientAddr, err)
			return
		}

		message := string(buffer[:n])
		log.Printf("Received from %s: %s", clientAddr, message)

		// Echo back to the client
		_, err = conn.Write(buffer[:n])
		if err != nil {
			log.Printf("Error writing to client %s: %v", clientAddr, err)
			return
		}
	}
}

func (s *MatchServer) addClient(addr string, conn net.Conn) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()
	s.clients[addr] = conn
}

func (s *MatchServer) removeClient(addr string) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()
	delete(s.clients, addr)
	log.Printf("Client disconnected: %s", addr)
}

func (s *MatchServer) Broadcast(message []byte) {
	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()

	for addr, conn := range s.clients {
		_, err := conn.Write(message)
		if err != nil {
			log.Printf("Error broadcasting to client %s: %v", addr, err)
		}
	}
}
