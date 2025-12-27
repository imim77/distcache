package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/imim77/distcache/cache"
	"github.com/imim77/distcache/proto"
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
	fmt.Println("connection made:", conn.RemoteAddr())
	defer func() {
		conn.Close()
	}()
	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			log.Println("parse command error:", err)
			break
		}
		fmt.Println(cmd)
		go s.handleCommand(conn, cmd)
	}
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:

	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	fmt.Printf("SET %s to %s\n", cmd.Key, cmd.Value)
	return s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL))

}
