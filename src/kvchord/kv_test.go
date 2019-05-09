package kvchord

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestKeysPerNode(t *testing.T) {
	fmt.Println("Test Keys Per Node")
	localhost := "127.0.0.1"
	initialServer := StartKVServer(localhost+":4000", "")
	s1 := StartKVServer(localhost+":5000", initialServer.chord.Ip)
	s2 := StartKVServer(localhost+":5002", initialServer.chord.Ip)
	s3 := StartKVServer(localhost+":5004", initialServer.chord.Ip)

	time.Sleep(3 * time.Second)
	initialServer.put(localhost+"10002", "HEHE")
	for index := 0; index < 10; index++ {
		initialServer.put(string(index), strconv.Itoa(index))
	}

	for index := 0; index < 10; index++ {
		val, _ := initialServer.get(string(index))
		fmt.Printf("%v: %v\n", index, val)
	}

	fmt.Println(initialServer.String())
	fmt.Println(s1.String())
	fmt.Println(s2.String())
	fmt.Println(s3.String())
}

func print(str string) {
	fmt.Println(str)
}
