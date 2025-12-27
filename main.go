package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/imim77/distcache/cache"
	"github.com/imim77/distcache/client"
	"github.com/imim77/distcache/proto"
)

func main() {
	listenAdr := flag.String("listenaddr", ":3000", "listen address of the server")
	leaderAddr := flag.String("leaderaddr", "", "listen address of the leader")
	flag.Parse()
	opts := ServerOpts{
		ListenAddr: *listenAdr,
		IsLeader:   len(*leaderAddr) == 0,
		LeaderAddr: *leaderAddr,
	}

	go func() {
		time.Sleep(time.Second * 2)
		for i := 0; i < 10; i++ {
			SendCommand()
			time.Sleep(time.Millisecond * 200)
		}

	}()
	s := NewServer(opts, cache.New())
	s.Start()
}

func SendCommand() {
	cmd := &proto.CommandSet{
		Key:   []byte("Foo"),
		Value: []byte("Bar"),
		TTL:   2,
	}

	client, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	conn.Write(cmd.Bytes())
}
