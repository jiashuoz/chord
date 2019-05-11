package main

import (
	"chord"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := chord.DefaultConfig()

	h, err := chord.MakeChord(cfg, "140.180.242.76:8001", "")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		time.Sleep(5 * time.Second)
		fmt.Println(h.String())
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-time.After(10 * time.Second)
	<-c
	// h.Stop()
}
