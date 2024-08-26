package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Acccept err:", err)
			continue
		}

		fmt.Println("new connection to the server:", conn.RemoteAddr())

		conn.Write([]byte("Sucessfully connected to Server:\n"))
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read err:", err)
			continue
		}

		s.msgch <- Message{
			from:    conn.LocalAddr().String(),
			payload: buf[:n],
		}

		conn.Write([]byte("Recevied your message!\n"))
	}
}

func main() {
	// Call the function
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgch {
			fmt.Printf("received msg from conn:%s || %s\n", msg.from, string(msg.payload))
		}
	}()
	log.Fatal(server.Start())
}
