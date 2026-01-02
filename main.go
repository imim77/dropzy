package main

import (
	"log"
	"time"

	"github.com/imim77/dropzy/server"
)

func main() {
	cfg := server.ServerConfig{
		Version:     "alpha",
		ListenAddr:  ":4000",
		GameVariant: server.TexasHoldem,
	}
	s := server.NewServer(cfg)
	go s.Start()

	time.Sleep(time.Second * 1)

	remoteCfg := server.ServerConfig{
		Version:     "alpha",
		ListenAddr:  ":5000",
		GameVariant: server.TexasHoldem,
	}
	remoteServer := server.NewServer(remoteCfg)

	go remoteServer.Start()

	if err := remoteServer.Connect(":4000"); err != nil {
		log.Fatal(err)
	}
	select {}

	//time.Sleep(time.Second * 1)
	//time.Sleep(time.Second)
	//remoteCfg := server.ServerConfig{
	//Version:    "GGPOKER v1.2 Alpha",
	//ListenAddr: ":2000",
	//}
	//remoteSrv := server.NewServer(remoteCfg)
	//go func() {
	//remoteSrv.Start()
	//}()
	//time.Sleep(time.Second * 1)
	//if err := remoteSrv.Connect(":4000"); err != nil {
	//fmt.Println(err)
	//}
	//select {}

}
