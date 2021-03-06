package chord

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTransferKey(t *testing.T) {
	ringsize := 8
	defaultConfig.ringSize = 3
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 1024)
	chords := make([]*ChordServer, ringsize)

	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	chords[1], _ = MakeChord(defaultConfig, testAddrs[1], testAddrs[0])
	chords[3], _ = MakeChord(defaultConfig, testAddrs[3], testAddrs[0])

	time.Sleep(2 * time.Second)
	chords[0].put(testAddrs[5], "world")
	chords[0].put(testAddrs[4], "world1")

	time.Sleep(1 * time.Second)
	fmt.Println(chords[0].String())
	chords[5], _ = MakeChord(defaultConfig, testAddrs[5], testAddrs[0])
	time.Sleep(1 * time.Second)
	chords[5].Leave()
	time.Sleep(1 * time.Second)
	for _, c := range chords {
		if c != nil {
			fmt.Println(c.String())
		}
	}
}

func TestLatency500Nodes(t *testing.T) {
	ringsize := 512
	defaultConfig.ringSize = 9 // this is not ringSize, this should be x, where 2^x = ringsize
	numOfNodes := 500
	jump := int(ringsize / numOfNodes)
	numOfLookup := 200

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	chords := make([]*ChordServer, ringsize)
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 1024) // ip addresses in order
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	success := 0
	var err error
	for index := 1; index < numOfNodes; index++ {
		chords[index*jump], err = MakeChord(defaultConfig, testAddrs[index*jump], testAddrs[0])
		if err == nil {
			success++
		}
	}

	fmt.Printf("successfully launched %d servers\n", success)

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println(chords[0].String())

	for i := 0; i < numOfLookup; i++ {
		chords[0].Lookup(fmt.Sprintf("%d:%f:%d", r1.Int(), r1.Float64(), r1.Int()))
		chords[0].tracerRWMu.Lock()
		fmt.Println(chords[0].tracer.Hops() + " " + chords[0].tracer.Latency())
		chords[0].tracerRWMu.Unlock()
	}
}

func TestLatency160Nodes(t *testing.T) {
	ringsize := 256
	defaultConfig.ringSize = 8
	numOfNodes := 160
	jump := int(ringsize / numOfNodes)
	numOfLookup := 200

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	chords := make([]*ChordServer, ringsize)
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000) // ip addresses in order
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	for index := 1; index < numOfNodes; index++ {
		chords[index*jump], _ = MakeChord(defaultConfig, testAddrs[index*jump], testAddrs[0])
	}

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println(chords[0].String())

	for i := 0; i < numOfLookup; i++ {
		chords[0].Lookup(fmt.Sprintf("%d:%f:%d", r1.Int(), r1.Float64(), r1.Int()))
		chords[0].tracerRWMu.Lock()
		fmt.Println(chords[0].tracer.Hops() + " " + chords[0].tracer.Latency())
		chords[0].tracerRWMu.Unlock()
	}
}

func TestLatency100Nodes(t *testing.T) {
	ringsize := 256
	defaultConfig.ringSize = 8
	numOfNodes := 100
	jump := int(ringsize / numOfNodes)
	numOfLookup := 200

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	chords := make([]*ChordServer, ringsize)
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000) // ip addresses in order
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	for index := 1; index < numOfNodes; index++ {
		chords[index*jump], _ = MakeChord(defaultConfig, testAddrs[index*jump], testAddrs[0])
	}

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println(chords[0].String())

	for i := 0; i < numOfLookup; i++ {
		chords[0].Lookup(fmt.Sprintf("%d:%f:%d", r1.Int(), r1.Float64(), r1.Int()))
		chords[0].tracerRWMu.Lock()
		fmt.Println(chords[0].tracer.Hops() + " " + chords[0].tracer.Latency())
		chords[0].tracerRWMu.Unlock()
	}
}

func TestLatency80Nodes(t *testing.T) {
	ringsize := 256
	defaultConfig.ringSize = 8
	numOfNodes := 80
	jump := int(ringsize / numOfNodes)
	numOfLookup := 200

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	chords := make([]*ChordServer, ringsize)
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000) // ip addresses in order
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	for index := 1; index < numOfNodes; index++ {
		chords[index*jump], _ = MakeChord(defaultConfig, testAddrs[index*jump], testAddrs[0])
	}

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println(chords[0].String())

	for i := 0; i < numOfLookup; i++ {
		chords[0].Lookup(fmt.Sprintf("%d:%f:%d", r1.Int(), r1.Float64(), r1.Int()))
		chords[0].tracerRWMu.Lock()
		fmt.Println(chords[0].tracer.Hops() + " " + chords[0].tracer.Latency())
		chords[0].tracerRWMu.Unlock()
	}
}

func TestLatency40Nodes(t *testing.T) {
	ringsize := 256
	defaultConfig.ringSize = 8
	numOfNodes := 40
	jump := int(ringsize / numOfNodes)
	numOfLookup := 200

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	chords := make([]*ChordServer, ringsize)
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000) // ip addresses in order
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	for index := 1; index < numOfNodes; index++ {
		chords[index*jump], _ = MakeChord(defaultConfig, testAddrs[index*jump], testAddrs[0])
	}

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println(chords[0].String())

	for i := 0; i < numOfLookup; i++ {
		chords[0].Lookup(fmt.Sprintf("%d:%f:%d", r1.Int(), r1.Float64(), r1.Int()))
		chords[0].tracerRWMu.Lock()
		fmt.Println(chords[0].tracer.Hops() + " " + chords[0].tracer.Latency())
		chords[0].tracerRWMu.Unlock()
	}
}

func TestLatency20Nodes(t *testing.T) {
	ringsize := 256
	defaultConfig.ringSize = 8
	numOfNodes := 20
	jump := int(ringsize / numOfNodes)
	numOfLookup := 200

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	chords := make([]*ChordServer, ringsize)
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000) // ip addresses in order
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	for index := 1; index < numOfNodes; index++ {
		chords[index*jump], _ = MakeChord(defaultConfig, testAddrs[index*jump], testAddrs[0])
	}

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println(chords[0].String())

	for i := 0; i < numOfLookup; i++ {
		chords[0].Lookup(fmt.Sprintf("%d:%f:%d", r1.Int(), r1.Float64(), r1.Int()))
		chords[0].tracerRWMu.Lock()
		fmt.Println(chords[0].tracer.Hops() + " " + chords[0].tracer.Latency())
		chords[0].tracerRWMu.Unlock()
	}
}

func TestLatency10Nodes(t *testing.T) {
	ringsize := 256
	defaultConfig.ringSize = 8 // 2^8 = 256
	numOfNodes := 10
	jump := int(ringsize / numOfNodes)
	numOfLookup := 200

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	chords := make([]*ChordServer, ringsize)
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000) // ip addresses in order
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	for index := 1; index < numOfNodes; index++ {
		chords[index*jump], _ = MakeChord(defaultConfig, testAddrs[index*jump], testAddrs[0])
	}

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println(chords[0].String())

	for i := 0; i < numOfLookup; i++ {
		chords[0].Lookup(fmt.Sprintf("%d:%f:%d", r1.Int(), r1.Float64(), r1.Int()))
		chords[0].tracerRWMu.Lock()
		fmt.Println(chords[0].tracer.Hops() + " " + chords[0].tracer.Latency())
		chords[0].tracerRWMu.Unlock()
	}
}

func TestWAN(t *testing.T) {
	defaultConfig.ringSize = 8
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000)
	chords := make([]*ChordServer, 100)
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	chords[1], _ = MakeChord(defaultConfig, testAddrs[1], testAddrs[0])

	time.Sleep(3 * time.Second)

	fmt.Println(chords[0].String())
	fmt.Println(chords[1].String())

	for _, c := range chords { // stop stabilize and fixfingers
		if c != nil {
			fmt.Println(c.String())
			c.StopFixFingers()
		}
	}

	chords[0].Lookup(testAddrs[1])
	fmt.Println(chords[0].tracer.String())
}

func TestLookup256Nodes(t *testing.T) {
	defaultConfig.ringSize = 8
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000)
	chords := make([]*ChordServer, 100)
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")

	for index := 1; index < 100; index++ {
		chords[index], _ = MakeChord(defaultConfig, testAddrs[index], chords[0].Ip)
	}

	time.Sleep(10 * time.Second)

	fmt.Println("Finger tables.....")
	for _, c := range chords {
		if c != nil {
			fmt.Println(c.String())
			c.StopFixFingers()
		}
	}

	fmt.Println()
	time.Sleep(5 * time.Second)

	fmt.Println("Start lookup")
	chords[0].Lookup(chords[99].Ip)
	fmt.Println(chords[0].tracer.String())
	fmt.Println("finish lookup")
}

func TestLookup32Nodes(t *testing.T) {
	numBits := 32
	defaultConfig.ringSize = 5
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000)
	chords := make([]*ChordServer, numBits)
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")

	for index := 1; index < numBits; index++ {
		chords[index], _ = MakeChord(defaultConfig, testAddrs[index], chords[0].Ip)
	}

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println("Finger tables.....")
	for _, c := range chords {
		if c != nil {
			fmt.Println(c.String())
			c.StopFixFingers()
		}
	}
	fmt.Println()
	time.Sleep(5 * time.Second)

	fmt.Println("Start lookup")
	chords[0].Lookup(chords[numBits-1].Ip)
	fmt.Println(chords[0].tracer.String())
	fmt.Println("finish lookup")
}

// This shows the max hop of a ring 2^4
func TestLookup16Nodes(t *testing.T) {
	defaultConfig.ringSize = 4
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000)
	chords := make([]*ChordServer, 16)
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")

	for index := 1; index < 16; index++ {
		chords[index], _ = MakeChord(defaultConfig, testAddrs[index], chords[0].Ip)
	}

	time.Sleep(5 * time.Second) // let it stabilize

	fmt.Println("Finger tables.....")
	for _, c := range chords {
		if c != nil {
			fmt.Println(c.String())
			c.StopFixFingers()
		}
	}
	fmt.Println()
	time.Sleep(5 * time.Second)

	fmt.Println("Start lookup")
	chords[0].Lookup(chords[14].Ip)
	fmt.Println(chords[0].tracer.String())
	fmt.Println("finish lookup")
}

func TestLookupCorrectness(t *testing.T) {

	defaultConfig.ringSize = 4
	testAddrs := reverseHash(defaultConfig.ringSize, "127.0.0.1", 5000)

	chords := make([]*ChordServer, 16)
	chords[0], _ = MakeChord(defaultConfig, testAddrs[0], "")
	chords[1], _ = MakeChord(defaultConfig, testAddrs[1], testAddrs[0])
	chords[3], _ = MakeChord(defaultConfig, testAddrs[3], testAddrs[0])
	chords[6], _ = MakeChord(defaultConfig, testAddrs[6], testAddrs[0])
	chords[7], _ = MakeChord(defaultConfig, testAddrs[7], testAddrs[0])
	chords[11], _ = MakeChord(defaultConfig, testAddrs[11], testAddrs[0])
	chords[14], _ = MakeChord(defaultConfig, testAddrs[14], testAddrs[0])

	time.Sleep(5 * time.Second)

	for {
		fmt.Println("Finger tables.....")
		for _, c := range chords {
			if c != nil {
				fmt.Println(c.String())
				c.StopFixFingers()
			}
		}
		fmt.Println()
		time.Sleep(5 * time.Second)
		break
	}

	fmt.Println("Start lookup")
	chords[0].Lookup(chords[14].Ip)
	fmt.Println(chords[0].tracer.String())
	fmt.Println("finish lookup")
}

// Test: node with conflicting ID is not allowed to join the ring
func TestConflictNodeID(t *testing.T) {
	numBits := 3
	fmt.Println("Test TestConflictNodeID")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, err := MakeChord(defaultConfig, testAddrs[0], "")
	chord1, err := MakeChord(defaultConfig, testAddrs[1], chord0.Node.Ip)

	time.Sleep(3 * time.Second)
	// Create a new chord instance with the same ID as chord1, needs a different ip address for rpc to work
	testAddrsDup := reverseHash(numBits, "127.0.0.1", 10000)
	_, err = MakeChord(defaultConfig, testAddrsDup[1], chord1.Node.Ip)

	if err != nil {
		DPrintf("TestStop %v", err)
	}
}

// Test graceful stopping chord instance
func TestStop(t *testing.T) {
	numBits := 3
	fmt.Println("Test TestStop")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)

	chord1, err := MakeChord(defaultConfig, testAddrs[1], "")
	checkError("TestStop", err)
	time.Sleep(3 * time.Second)

	fmt.Println(chord1.String())

	chord1.Stop()
	time.Sleep(3 * time.Second)
}
