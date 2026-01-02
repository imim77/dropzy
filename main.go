package main

import (
	"fmt"
	"time"

	"github.com/imim77/dropzy/server"
)

func main() {
	cfg := server.ServerConfig{
		Version:    "GGPOKER v1.2 Alpha",
		ListenAddr: ":4000",
	}
	s := server.NewServer(cfg)

	go func() {
		s.Start()
	}()
	time.Sleep(time.Second * 1)
	time.Sleep(time.Second)
	remoteCfg := server.ServerConfig{
		Version:    "GGPOKER v1.2 Alpha",
		ListenAddr: ":2000",
	}
	remoteSrv := server.NewServer(remoteCfg)
	go func() {
		remoteSrv.Start()
	}()
	if err := remoteSrv.Connect(":4000"); err != nil {
		fmt.Println(err)
	}

}
