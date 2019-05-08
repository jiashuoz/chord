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
	chord0, err := MakeChord(testAddrs[0], nil)
	chord1, err := MakeChord(testAddrs[1], chord0.Node)

	time.Sleep(3 * time.Second)
	// Create a new chord instance with the same ID as chord1, needs a different ip address for rpc to work
	testAddrsDup := reverseHash(numBits, "127.0.0.1", 10000)
	_, err = MakeChord(testAddrsDup[1], chord1.Node)

	if err != nil {
		DPrintf("TestStop %v", err)
	}
}

// Test graceful stopping chord instance
func TestStop(t *testing.T) {
	fmt.Println("Test TestStop")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)

	chord1, err := MakeChord(testAddrs[1], nil)
	checkError("TestStop", err)
	time.Sleep(3 * time.Second)

	fmt.Println(chord1.String())

	chord1.Stop("force stop chord1")
	time.Sleep(3 * time.Second)
}
