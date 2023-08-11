package client

// Client for dbserver/slowdb

import (
	"fmt"
	"github.com/capotej/groupcache-db-experiment/api"
	"github.com/capotej/groupcache-db-experiment/rpc"
)

type Client struct{}

func (c *Client) Get(key string) string {
	client, err := rpc.DialHTTP("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("error %s\n", err)
		return ""
	}

	fmt.Printf("client dail local db: Get %+v\n", key)
	args := &api.Load{key}

	var reply api.ValueResult
	err = client.Call("Server.Get", args, &reply)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return ""
	}
	return string(reply.Value)
}

func (c *Client) Set(key string, value string) {
	client, err := rpc.DialHTTP("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("error %s\n", err)
		return
	}

	fmt.Printf("client dail local db: Set %+v:%+v\n", key, value)
	args := &api.Store{key, value}

	var reply int
	err = client.Call("Server.Set", args, &reply)
	if err != nil {
		fmt.Printf("error %s\n", err)
		return
	}
}
