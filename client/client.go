package client

import (
	"context"
	"net"

	"github.com/imim77/distcache/proto"
)

type Options struct{}

type Client struct {
	conn net.Conn
}

func New(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Set(ctx context.Context, key, value []byte) (any, error) {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
	}
	return nil, nil
}
