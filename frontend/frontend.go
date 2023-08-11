package main

// This represents a cache front end server, that front slowdb/slowserver requests

import (
	"context"
	"flag"
	"fmt"
	"github.com/capotej/groupcache-db-experiment/api"
	"github.com/capotej/groupcache-db-experiment/client"
	"github.com/capotej/groupcache-db-experiment/groupcache"
	"github.com/capotej/groupcache-db-experiment/rpc"
	"net"
	"net/http"
	"os"
	"strconv"
)

type Frontend struct {
	cacheGroup *groupcache.Group
}

func (s *Frontend) Get(args *api.Load, reply *api.ValueResult) error {
	var data []byte
	fmt.Printf("cli asked for %s from groupcache\n", args.Key)
	err := s.cacheGroup.Get(context.Background(), args.Key,
		groupcache.AllocatingByteSliceSink(&data))

	reply.Value = string(data)
	return err
}

func NewServer(cacheGroup *groupcache.Group) *Frontend {
	server := new(Frontend)
	server.cacheGroup = cacheGroup
	return server
}

func (s *Frontend) Start(port string) {

	rpc.Register(s)

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", port)
	if e != nil {
		fmt.Println("fatal")
	}

	http.Serve(l, nil)
}

func main() {
	var port = flag.String("port", "8001", "groupcache port")
	flag.Parse()

	peers := groupcache.NewHTTPPool("http://localhost:" + *port)
	peers.Set("http://localhost:8001", "http://localhost:8002", "http://localhost:8003")

	client := new(client.Client)
	var stringcache = groupcache.NewGroup("SlowDBCache", 64<<20, groupcache.GetterFunc(
		// 执行s.cacheGroup.Get
		// lookupCache
		// g.getFromPeer
		// g.getLocally 最后一步，就是执行下面的函数
		// result 通过 groupcache.AllocatingByteSliceSink(&data) 的 SetBytes 方法，set 到data变量，data通常是一个引用
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			result := client.Get(key)
			fmt.Printf("GetterFunc asking for %s from dbserver\n", key)
			dest.SetBytes([]byte(result))
			return nil
		}))
	frontendServer := NewServer(stringcache)

	i, err := strconv.Atoi(*port)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}

	// 9001 9002 9003
	var frontEndport = ":" + strconv.Itoa(i+1000)
	go frontendServer.Start(frontEndport)

	fmt.Println(stringcache)
	fmt.Println("cachegroup slave starting on " + *port)
	fmt.Println("frontend starting on " + frontEndport)
	http.ListenAndServe("127.0.0.1:"+*port, http.HandlerFunc(peers.ServeHTTP))
}
