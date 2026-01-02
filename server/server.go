package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

const (
	TexasHoldem GameVariant = iota
	Other
)

func (gv GameVariant) String() string {
	switch gv {
	case TexasHoldem:
		return "TEXAS_HOLDEM"
	case Other:
		return "OTHER"
	default:
		return "unknown"
	}
}

type ServerConfig struct {
	Version     string
	ListenAddr  string
	GameVariant GameVariant
}

type Server struct {
	ServerConfig

	transport *TCPTransport
	mu        *sync.RWMutex
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	msgCh     chan *Message
	delPeer   chan *Peer // kanal za brisanje Peerova(konekcija)
}

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
		delPeer:      make(chan *Peer),
	}
	tr := NewTCPTransport(s.ListenAddr)
	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer
	s.transport = tr
	return s
}

func (s *Server) Start() {
	go s.loop()

	logrus.WithFields(logrus.Fields{"port": s.ListenAddr, "type": "Texas hold em"}).Info("Started new server")
	s.transport.ListenAndAccept()

}

func (s *Server) SendHandshake(p *Peer) error {
	hs := &Handshake{
		GameVariant: s.GameVariant,
		Version:     s.Version,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}
	return p.Send(buf.Bytes())
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player disconnected")
			addr := peer.conn.RemoteAddr()
			delete(s.peers, addr)

		case peer := <-s.addPeer:
			s.SendHandshake(peer)
			if err := s.handshake(peer); err != nil {
				logrus.Errorf("handshake with incoming peer failed: %s", err)
				continue
			} // prvo provjerimo handshake da vidimo da li se peer može spojiti

			go peer.ReadLoop(s.msgCh)

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("handshake successfull, new device connected")
			s.peers[peer.conn.RemoteAddr()] = peer

		case msg := <-s.msgCh:

			if err := s.handleMessage(msg); err != nil {
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

type Handshake struct {
	Version     string
	GameVariant GameVariant
}

// ako se peer-ovi(konekcije) "ne slože" na handshakeu onda dropamo peer koji se pokusava spojiti
func (s *Server) handshake(p *Peer) error {
	hs := &Handshake{}
	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return err
	}

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("invalid GameVarient: %s", hs.GameVariant)
	}
	if s.Version != hs.Version {
		return fmt.Errorf("invalid Version: %s", hs.Version)
	}
	logrus.WithFields(logrus.Fields{"peer": p.conn.RemoteAddr(), "version": hs.Version, "variant": hs.GameVariant}).Info("recieved handshake")
	return nil
}

func (s *Server) handleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}
