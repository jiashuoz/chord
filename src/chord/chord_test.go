package chord

import (
	"fmt"
	"testing"
	"time"
)

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
		for _, c := range chords {
			if c != nil {
				fmt.Println(c.String())
			}
		}
		time.Sleep(5 * time.Second)
	}

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
