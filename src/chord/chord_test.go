package chord

import (
	"fmt"
	"testing"
	"time"
)

// Test: node with conflicting ID is not allowed to join the ring
func TestConflictNodeID(t *testing.T) {
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
	fmt.Println("Test TestStop")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)

	chord1, err := MakeChord(defaultConfig, testAddrs[1], "")
	checkError("TestStop", err)
	time.Sleep(3 * time.Second)

	fmt.Println(chord1.String())

	chord1.Stop()
	time.Sleep(3 * time.Second)
}
