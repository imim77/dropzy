package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

type Message struct {
	Payload io.Reader
	From    net.Addr
}

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
func (p *Peer) ReadLoop(msgch chan *Message) {
	defer p.conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}
		msgch <- &Message{
			Payload: bytes.NewReader(buf[:n]),
			From:    p.conn.RemoteAddr(),
		}

	}

}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	AddPeer    chan *Peer
	DelPeer    chan *Peer
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: addr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		peer := &Peer{ // ovdje kreiramo novi peer
			conn: conn,
		}
		t.AddPeer <- peer

	}

	return fmt.Errorf("TCP transport stopped reason: ?")
}
