package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/imim77/distcache/cache"
	"github.com/imim77/distcache/client"
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
		client, err := client.New(":3000", client.Options{})
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < 10; i++ {
			SendCommand(client)
			time.Sleep(time.Millisecond * 200)
		}
		client.Close()
		time.Sleep(time.Second * 1)

	}()
	s := NewServer(opts, cache.New())
	s.Start()
}

func SendCommand(c *client.Client) {

	_, err := c.Set(context.Background(), []byte("gg"), []byte("Anthony"), 0)
	if err != nil {
		log.Fatal(err)
	}

}
