package main

import (
	"fmt"
	"io"
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
	//fmt.Println("connection made:", conn.RemoteAddr())
	defer func() {
		conn.Close()
	}()
	for {
		cmd, err := proto.ParseCommand(conn)
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("parse command error:", err)
			break
		}
		fmt.Println(cmd)
		go s.handleCommand(conn, cmd)
	}

	//fmt.Println("connection closed: ", conn.RemoteAddr())
}

func (s *Server) handleCommand(conn net.Conn, cmd any) {
	switch v := cmd.(type) {
	case *proto.CommandSet:
		s.handleSetCommand(conn, v)
	case *proto.CommandGet:
		s.handleGetCommand(conn, v)

	}
}

func (s *Server) handleSetCommand(conn net.Conn, cmd *proto.CommandSet) error {
	//fmt.Printf("SET %s to %s\n", cmd.Key, cmd.Value)
	resp := proto.ResponseSet{}
	if err := s.cache.Set(cmd.Key, cmd.Value, time.Duration(cmd.TTL)); err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}

	resp.Status = proto.StatusOK
	_, err := conn.Write(resp.Bytes())
	return err

}

func (s *Server) handleGetCommand(conn net.Conn, cmd *proto.CommandGet) error {
	//log.Printf("GET %s", cmd.Key)
	resp := proto.ResponseGet{}
	value, err := s.cache.Get(cmd.Key)
	if err != nil {
		resp.Status = proto.StatusError
		_, err := conn.Write(resp.Bytes())
		return err
	}
	resp.Status = proto.StatusOK
	resp.Value = value
	_, err = conn.Write(resp.Bytes())

	return nil
}
