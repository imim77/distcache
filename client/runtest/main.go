package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/imim77/distcache/client"
)

func main() {
	SendStuff()
}

func SendStuff() {

	client, err := client.New(":3000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {

		var (
			key   = []byte(fmt.Sprintf("key_%d", i))
			value = []byte(fmt.Sprintf("val_%d", i))
		)

		err = client.Set(context.Background(), key, value, 0)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Second)
	}
	client.Close()

}
