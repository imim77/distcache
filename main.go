package main

import (
	"flag"

	"github.com/imim77/distcache/cache"
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
