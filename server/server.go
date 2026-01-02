package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"
)

// Peer is a connection
type Peer struct {
	conn net.Conn
}

func (p *Peer) Send(b []byte) error {
	if _, err := p.conn.Write(b); err != nil {
		return err
	}
	return nil

}

type ServerConfig struct {
	Version    string
	ListenAddr string
}

type Message struct {
	Payload io.Reader
	From    net.Addr
}

type Server struct {
	ServerConfig
	handler  Handler
	listener net.Listener
	mu       *sync.RWMutex
	peers    map[net.Addr]*Peer
	addPeer  chan *Peer
	msgCh    chan *Message
	delPeer  chan *Peer // kanal za brisanje Peerova(konekcija)
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
		delPeer:      make(chan *Peer),
		handler:      &DefaultHandler{},
	}
}

func (s *Server) Start() {
	go s.loop()
	if err := s.listen(); err != nil {
		panic(err)
	}
	fmt.Printf("[INFO] Server started on port: %s\n", s.ListenAddr)
	s.acceptLoop()

}

func (s *Server) listen() error {
	ln, err := net.Listen("tcp", s.ListenAddr)

	if err != nil {
		return err
	}
	s.listener = ln
	return nil

}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			addr := peer.conn.RemoteAddr()
			delete(s.peers, addr)
			fmt.Printf("payer disconnected %s\n", peer.conn.RemoteAddr())
		case peer := <-s.addPeer:
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("new device conneted %s\n", peer.conn.RemoteAddr())
		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}
	s.addPeer <- peer
	return peer.Send([]byte(s.Version))

}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		peer := &Peer{ // Ovjde kreiramo peer
			conn: conn,
		}

		s.addPeer <- peer

		peer.Send([]byte(s.Version))
		go s.handleConn(peer)
	}

}

func (s *Server) handleConn(p *Peer) {
	defer func() { s.delPeer <- p }()

	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}

		s.msgCh <- &Message{
			Payload: bytes.NewReader(buf[:n]),
			From:    p.conn.RemoteAddr(),
		}
	}

}
