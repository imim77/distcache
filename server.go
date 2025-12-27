package main

import (
	"context"
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
	followers map[net.Conn]struct{}
	cache     cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		followers:  make(map[net.Conn]struct{}),
	}

}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}
	log.Printf("Server starting on port [%s]\n", s.ListenAddr)

	if !s.IsLeader {
		go func() {
			conn, err := net.Dial("tcp", s.LeaderAddr)
			fmt.Println("connected with leader: ", s.LeaderAddr)
			if err != nil {
				log.Fatal(err)
			}
			s.handleConn(conn)
		}()
	}

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
	msg, err := parseMessage(rawCMD)
	if err != nil {
		fmt.Println("failed to parse command", err)
		conn.Write([]byte(err.Error()))
		return
	}
	fmt.Printf("recieved command %s\n", msg.Cmd)
	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCmd(conn, msg)
	case CMDGet:
		err = s.handleGetCmd(conn, msg)
	}

	if err != nil {
		fmt.Println("failed to handle command", err)
		conn.Write([]byte(err.Error()))
	}

}

func (s *Server) handleGetCmd(conn net.Conn, msg *Message) error {
	val, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}
	_, err = conn.Write(val)
	return err

}
func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	if err := s.cache.Set(msg.Key, msg.Value, msg.TTL); err != nil {
		return err
	}

	go s.sendToFollowers(context.TODO(), msg)

	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	for conn := range s.followers {
		fmt.Println("forwarding key to follower")
		_, err := conn.Write(msg.ToBytes())
		rawMsg := msg.ToBytes()
		fmt.Println("fowarding rawmsg to follower: ", string(rawMsg))
		if err != nil {
			fmt.Println("write to follower error", err)
			continue
		}
	}
	return nil
}
