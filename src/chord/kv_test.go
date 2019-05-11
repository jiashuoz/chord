package chord

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestLookupLatency(t *testing.T) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	serverCount := 3

	randomAddr := ipGenerator("127.0.0.1", serverCount) // generates 100 unique random ip address
	// inorderAddr := reverseHash(config.ringSize, "127.0.0.1", 5000)

	initialServer, err := MakeChord(defaultConfig, randomAddr[0], "")

	servers := make([]*ChordServer, serverCount)
	for i := 1; i < serverCount; i++ {
		servers[i], err = MakeChord(defaultConfig, randomAddr[i], initialServer.Node.Ip)
		checkErrorGrace("", err)
	}

	fmt.Println("finished launching servers")

	time.Sleep(40 * time.Second)

	for i := 1; i < 2; i++ {
		fmt.Println(servers[i].String())
	}

	for i := 0; i < 2000; i++ {
		initialServer.Lookup(strconv.Itoa(r1.Int()))
		initialServer.tracerRWMu.Lock()
		fmt.Println(initialServer.tracer.Hops() + " " + initialServer.tracer.Latency())
		initialServer.tracerRWMu.Unlock()
	}

	time.Sleep(1 * time.Second)

}

func TestKeysPerNodeAdv(t *testing.T) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	serverCount := 100

	addrs := ipGenerator("127.0.0.1", serverCount) // generates 100 unique random ip address
	initialServer, err := MakeChord(defaultConfig, addrs[0], "")
	checkErrorGrace("", err)

	keyCount := 5000

	fmt.Printf("Launching %v servers, loading %v keys\n", serverCount, keyCount)

	servers := make([]*ChordServer, serverCount)
	for i := 1; i < serverCount; i++ {
		servers[i], err = MakeChord(defaultConfig, addrs[i], initialServer.Node.Ip)
		checkErrorGrace("", err)
	}

	fmt.Println("finished launching servers")

	time.Sleep(10 * time.Second)

	for i := 0; i < keyCount; i++ {
		initialServer.put(strconv.Itoa(r1.Int()), strconv.Itoa(i))
	}

	fmt.Println("finished putting keys")

	fmt.Println("printing results...")
	totalKeys := initialServer.keyCount()
	fmt.Printf("number of keys: %v\n", totalKeys)
	for i := 1; i < serverCount; i++ {
		eachCount := servers[i].keyCount()
		totalKeys += eachCount
		fmt.Printf("number of keys: %v\n", eachCount)
	}
	if totalKeys != keyCount {
		fmt.Fprintf(os.Stderr, "totalKeys %v != keyCount %v", totalKeys, keyCount)
	}
	fmt.Printf("total key: %v\n", totalKeys)
}
