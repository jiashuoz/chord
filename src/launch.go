package main

import (
	"chord"
	"github.com/jiashuoz/chord/chordrpc"
	"log"
	"os"
	"os/signal"
	"time"
)

func createNode(ip string, joinNode *chordrpc.Node) (*chord.ChordServer, error) {

	n, err := chord.MakeChord(ip, joinNode)
	return n, err
}

func main() {

	h, err := createNode("0.0.0.0:8001", nil)
	if err != nil {
		log.Fatalln(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-time.After(10 * time.Second)
	<-c
	h.Stop()

}
