package main

import (
	"context"
	"flag"
	"fmt"
	"log"

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

	s := NewServer(opts, cache.New())
	s.Start()
}

func SendStuff() {

	for i := 0; i < 100; i++ {
		go func(i int) {
			client, err := client.New(":3000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}
			var (
				key   = []byte(fmt.Sprintf("key_%d", i))
				value = []byte(fmt.Sprintf("val_%d", i))
			)

			err = client.Set(context.Background(), key, value, 0)
			if err != nil {
				log.Fatal(err)
			}
			resp, err := client.Get(context.Background(), key)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(resp))

			client.Close()
		}(i)
	}
}
