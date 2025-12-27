package main

import (
	"fmt"
	"log"
	"net"

	"github.com/imim77/distcache/cache"
)

type ServerOpts struct {
	ListenAddr string
	LeaderAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}

}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}
	log.Printf("Server starting on port [%s]\n", s.ListenAddr)

	for {
		con, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue
		}
		go s.handleConn(con)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	buff := make([]byte, 2048)

	fmt.Println("connection made:", conn.RemoteAddr())
	defer func() {
		conn.Close()
	}()
	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Printf("conn read error: %s\n", err)
			break
		}

		msg := buff[:n]
		s.handleCommand(conn, msg)
		//fmt.Println(string(msg))
	}
}

func (s *Server) handleCommand(conn net.Conn, rawCMD []byte) {
}
