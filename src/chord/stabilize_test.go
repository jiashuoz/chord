package chord

import (
	"bytes"
	"fmt"
	"github.com/jiashuoz/chord/chordrpc"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestConcurrentJoin2(t *testing.T) {
	fmt.Println("Test TestConcurrentJoin 2")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)

	chord0, _ := MakeChord(defaultConfig, testAddrs[0], "")

	launchChord := func(index int, joinNode *chordrpc.Node) {
		c, _ := MakeChord(defaultConfig, testAddrs[index], joinNode.Ip)
		fmt.Printf("initializing %v %v %v\n", index, testAddrs[index], c.Id)
		time.Sleep(5 * time.Second)
		fmt.Println(c.String())
	}

	go launchChord(1, chord0.Node)
	go launchChord(3, chord0.Node)
	go launchChord(6, chord0.Node)
	go launchChord(7, chord0.Node)
	go launchChord(11, chord0.Node)
	go launchChord(14, chord0.Node)

	time.Sleep(5 * time.Second)

	for index := 0; index < 1000; index++ {
		fmt.Println(chord0.Lookup(strconv.Itoa(index)))
	}

	fmt.Println(chord0.String())

	time.Sleep(5 * time.Second)
}

// Concurrent join, 1 3 6
func TestConcurrentJoin1(t *testing.T) {
	fmt.Println("Test TestConcurrentJoin 1")

	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, _ := MakeChord(defaultConfig, testAddrs[0], "")
	go func(index int) {
		fmt.Printf("initializing %v\n", index)
		c1, _ := MakeChord(defaultConfig, testAddrs[index], chord0.Node.Ip)
		time.Sleep(5 * time.Second)
		fmt.Println(c1.String())
	}(1)

	go func(index int) {
		fmt.Printf("initializing %v\n", index)
		c6, _ := MakeChord(defaultConfig, testAddrs[index], chord0.Node.Ip)
		time.Sleep(5 * time.Second)
		fmt.Println(c6.String())
	}(6)

	go func(index int) {
		fmt.Printf("initializing %v\n", index)
		c3, _ := MakeChord(defaultConfig, testAddrs[index], chord0.Node.Ip)
		time.Sleep(5 * time.Second)
		fmt.Println(c3.String())
	}(3)

	time.Sleep(5 * time.Second)
	fmt.Println(chord0.String())
	// Wait for all goroutines to print
	time.Sleep(10 * time.Second)
}

// Concurrent join 7 nodes
func TestConcurrentJoin0(t *testing.T) {
	fmt.Println("Test TestConcurrentJoin 0")

	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, err := MakeChord(defaultConfig, testAddrs[0], "")
	for index := 1; index <= 7; index++ {
		go func(index int) {
			c, _ := MakeChord(defaultConfig, testAddrs[index], chord0.Node.Ip)
			time.Sleep(5 * time.Second)
			fmt.Println(c.String())
		}(index)
	}

	if err != nil {
		checkError("TestStabilize2", err)
	}

	time.Sleep(5 * time.Second)
	fmt.Println(chord0.String())
	// Wait for all goroutines to print
	time.Sleep(10 * time.Second)
}

// Ring in figure 5, with 8 nodes
func TestStabilize2(t *testing.T) {
	fmt.Println("Test Stabilize2")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, err := MakeChord(defaultConfig, testAddrs[0], "")
	chord1, err := MakeChord(defaultConfig, testAddrs[1], chord0.Node.Ip)
	chord2, err := MakeChord(defaultConfig, testAddrs[2], chord1.Node.Ip)
	chord3, err := MakeChord(defaultConfig, testAddrs[3], chord2.Node.Ip)
	chord4, err := MakeChord(defaultConfig, testAddrs[4], chord3.Node.Ip)
	chord5, err := MakeChord(defaultConfig, testAddrs[5], chord2.Node.Ip)
	chord6, err := MakeChord(defaultConfig, testAddrs[6], chord1.Node.Ip)
	chord7, err := MakeChord(defaultConfig, testAddrs[7], chord0.Node.Ip)

	if err != nil {
		checkError("TestStabilize2", err)
	}

	for {
		fmt.Println(chord0.String())
		fmt.Println(chord1.String())
		fmt.Println(chord2.String())
		fmt.Println(chord3.String())
		fmt.Println(chord4.String())
		fmt.Println(chord5.String())
		fmt.Println(chord6.String())
		fmt.Println(chord7.String())
		time.Sleep(10 * time.Second)
	}
}

// Figure 3 setup, don't change
func TestStabilize0(t *testing.T) {
	fmt.Println("Test Stabilize0, figure3")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, err := MakeChord(defaultConfig, testAddrs[0], "")
	checkError("TestStabilize0", err)
	chord1, err := MakeChord(defaultConfig, testAddrs[1], chord0.Node.Ip)
	checkError("TestStabilize0", err)
	chord3, err := MakeChord(defaultConfig, testAddrs[3], chord0.Node.Ip)
	checkError("TestStabilize0", err)

	for {
		fmt.Println(chord0.String())
		fmt.Println(chord1.String())
		fmt.Println(chord3.String())
		time.Sleep(5 * time.Second)
	}
}

// Figure 5 setup in the paper, don't change
func TestStabilize1(t *testing.T) {
	fmt.Println("Test Stabilize1, figure5")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, err := MakeChord(defaultConfig, testAddrs[0], "")
	checkError("TestStabilize1", err)
	chord1, err := MakeChord(defaultConfig, testAddrs[1], chord0.Node.Ip)
	checkError("TestStabilize1", err)
	chord3, err := MakeChord(defaultConfig, testAddrs[3], chord0.Node.Ip)
	checkError("TestStabilize1", err)
	chord6, err := MakeChord(defaultConfig, testAddrs[6], chord0.Node.Ip)
	checkError("TestStabilize1", err)

	for {
		fmt.Println(chord0.String())
		fmt.Println(chord1.String())
		fmt.Println(chord3.String())
		fmt.Println(chord6.String())
		time.Sleep(5 * time.Second)
	}
}

func TestJoin(t *testing.T) {
	fmt.Println("Test Join")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, err := MakeChord(defaultConfig, testAddrs[0], "")
	checkError("TestJoin", err)

	for {
		fmt.Println(chord0.String())
		fmt.Println()
		time.Sleep(1 * time.Second)
	}
}

func TestMake(t *testing.T) {
	fmt.Println("Test MakeChord")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, err := MakeChord(defaultConfig, testAddrs[0], "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println(chord0.String())

	pred, err := chord0.getPredecessorRPC(chord0.Node)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(pred.String())

	succ, err := chord0.getSuccessorRPC(chord0.Node)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(succ.String())

	closest, err := chord0.findClosestPrecedingNodeRPC(chord0.Node, intToByteArray(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(closest.String())

	result, err := chord0.findSuccessorRPC(chord0.Node, intToByteArray(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	fmt.Println(result.String())
}

func TestFindSuccessor(t *testing.T) {
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, _ := MakeChord(defaultConfig, testAddrs[0], "") // id 0
	chord1, _ := MakeChord(defaultConfig, testAddrs[1], "") // id 1
	chord3, _ := MakeChord(defaultConfig, testAddrs[3], "") // id 3

	node0 := &chordrpc.Node{Id: intToByteArray(0), Ip: testAddrs[0]}
	node1 := &chordrpc.Node{Id: intToByteArray(1), Ip: testAddrs[1]}
	node3 := &chordrpc.Node{Id: intToByteArray(3), Ip: testAddrs[3]}

	chord0.fingerTable[0] = node1
	chord0.fingerTable[1] = node3
	chord0.fingerTable[2] = node0
	chord0.predecessor = node3

	chord1.fingerTable[0] = node3
	chord1.fingerTable[1] = node3
	chord1.fingerTable[2] = node0
	chord1.predecessor = node0

	chord3.fingerTable[0] = node0
	chord3.fingerTable[1] = node0
	chord3.fingerTable[2] = node0
	chord3.predecessor = node1

	// ba = byte array
	ba0 := intToByteArray(0)
	ba1 := intToByteArray(1)
	ba2 := intToByteArray(2)
	ba3 := intToByteArray(3)
	ba4 := intToByteArray(4)
	ba5 := intToByteArray(5)
	ba6 := intToByteArray(6)
	ba7 := intToByteArray(7)

	// FindSuccessor of key 0, call three servers, result should always be 0
	key0Succ, err := chord0.findSuccessor(ba0)
	checkError("TestFindSuccessor", err)
	if !bytes.Equal(key0Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key0Succ.Id)
	}

	key0Succ, err = chord1.findSuccessor(ba0)
	if !bytes.Equal(key0Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key0Succ.Id)
	}

	key0Succ, err = chord3.findSuccessor(ba0)
	if !bytes.Equal(key0Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key0Succ.Id)
	}

	// FindSuccessor of key 1, call three servers, result should always be 1
	key1Succ, err := chord0.findSuccessor(ba1)
	if !bytes.Equal(key1Succ.Id, intToByteArray(1)) {
		t.Errorf("Find successor from node %d got = %d; want 1", chord0.GetId(), key1Succ.Id)
	}

	key1Succ, err = chord1.findSuccessor(ba1)
	if !bytes.Equal(key1Succ.Id, intToByteArray(1)) {
		t.Errorf("Find successor from node %d got = %d; want 1", chord1.GetId(), key1Succ.Id)
	}

	key1Succ, err = chord3.findSuccessor(ba1)
	if !bytes.Equal(key1Succ.Id, intToByteArray(1)) {
		t.Errorf("Find successor from node %d got = %d; want 1", chord3.GetId(), key1Succ.Id)
	}

	// FindSuccessor of key 2, call three servers, result should always be 3
	key2Succ, err := chord0.findSuccessor(ba2)
	if !bytes.Equal(key2Succ.Id, intToByteArray(3)) {
		t.Errorf("Find successor from node %d got = %d; want 3", chord0.GetId(), key2Succ.Id)
	}

	key2Succ, err = chord1.findSuccessor(ba2)
	if !bytes.Equal(key2Succ.Id, intToByteArray(3)) {
		t.Errorf("Find successor from node %d got = %d; want 3", chord1.GetId(), key2Succ.Id)
	}

	key2Succ, err = chord3.findSuccessor(ba2)
	if !bytes.Equal(key2Succ.Id, intToByteArray(3)) {
		t.Errorf("Find successor from node %d got = %d; want 3", chord3.GetId(), key2Succ.Id)
	}

	// FindSuccessor of key 3, call three servers, result should always be 3
	key3Succ, err := chord0.findSuccessor(ba3)
	if !bytes.Equal(key3Succ.Id, intToByteArray(3)) {
		t.Errorf("Find successor from node %d got = %d; want 3", chord0.GetId(), key3Succ.Id)
	}

	key2Succ, err = chord1.findSuccessor(ba3)
	if !bytes.Equal(key3Succ.Id, intToByteArray(3)) {
		t.Errorf("Find successor from node %d got = %d; want 3", chord1.GetId(), key3Succ.Id)
	}

	key2Succ, err = chord3.findSuccessor(ba3)
	if !bytes.Equal(key3Succ.Id, intToByteArray(3)) {
		t.Errorf("Find successor from node %d got = %d; want 3", chord3.GetId(), key3Succ.Id)
	}

	// FindSuccessor of key 4, call three servers, result should always be 0
	key4Succ, err := chord0.findSuccessor(ba4)
	if !bytes.Equal(key4Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key4Succ.Id)
	}

	key4Succ, err = chord1.findSuccessor(ba4)
	if !bytes.Equal(key4Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key4Succ.Id)
	}

	key4Succ, err = chord3.findSuccessor(ba4)
	if !bytes.Equal(key4Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key4Succ.Id)
	}

	// FindSuccessor of key 5, call three servers, result should always be 0
	key5Succ, err := chord0.findSuccessor(ba5)
	if !bytes.Equal(key5Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key5Succ.Id)
	}

	key5Succ, err = chord1.findSuccessor(ba5)
	if !bytes.Equal(key5Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key5Succ.Id)
	}

	key5Succ, err = chord3.findSuccessor(ba5)
	if !bytes.Equal(key5Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key5Succ.Id)
	}

	// FindSuccessor of key 6, call three servers, result should always be 0
	key6Succ, err := chord0.findSuccessor(ba6)
	if !bytes.Equal(key6Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key6Succ.Id)
	}

	key6Succ, err = chord1.findSuccessor(ba6)
	if !bytes.Equal(key6Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key6Succ.Id)
	}

	key6Succ, err = chord3.findSuccessor(ba6)
	if !bytes.Equal(key6Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key6Succ.Id)
	}

	// FindSuccessor of key 7, call three servers, result should always be 0
	key7Succ, err := chord0.findSuccessor(ba7)
	if !bytes.Equal(key7Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key7Succ.Id)
	}

	key7Succ, err = chord1.findSuccessor(ba7)
	if !bytes.Equal(key7Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key7Succ.Id)
	}

	key7Succ, err = chord3.findSuccessor(ba7)
	if !bytes.Equal(key7Succ.Id, intToByteArray(0)) {
		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key6Succ.Id)
	}
}

// func TestGetSuccessor(t *testing.T) {
// 	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
// 	chord0 := MakeServer(testAddrs[0]) // id 0
// 	chord1 := MakeServer(testAddrs[1]) // id 1
// 	server2 := MakeServer(testAddrs[2]) // id 2
// 	chord3 := MakeServer(testAddrs[3]) // id 3

// 	node0 := MakeNode(testAddrs[0])
// 	node1 := MakeNode(testAddrs[1])
// 	node2 := MakeNode(testAddrs[2])
// 	node3 := MakeNode(testAddrs[3])

// 	chord0.fingerTable[0] = node1
// 	chord0.fingerTable[1] = node2
// 	chord0.fingerTable[2] = node0

// 	chord1.fingerTable[0] = node2
// 	chord1.fingerTable[1] = node3
// 	chord1.fingerTable[2] = node0

// 	server2.fingerTable[0] = node3
// 	server2.fingerTable[1] = node0
// 	server2.fingerTable[2] = node0

// 	chord3.fingerTable[0] = node0
// 	chord3.fingerTable[1] = node0
// 	chord3.fingerTable[2] = node0

// 	// run(chord0)
// 	// run(server2)
// 	// run(chord1)
// 	// run(chord3)

// 	fmt.Println(chord0.String(true))
// 	fmt.Println(chord1.String(true))
// 	fmt.Println(server2.String(true))
// 	fmt.Println(chord3.String(true))

// 	// ba = byte array
// 	ba0 := intToByteArray(0)
// 	ba1 := intToByteArray(1)
// 	ba2 := intToByteArray(2)
// 	ba3 := intToByteArray(3)
// 	ba4 := intToByteArray(4)
// 	ba5 := intToByteArray(5)
// 	ba6 := intToByteArray(6)
// 	ba7 := intToByteArray(7)

// 	// FindSuccessor of key 0, call three servers, result should always be 0
// 	key0Succ, err := chord0.findSuccessor(ba0)
// 	if !bytes.Equal(key0Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key0Succ.Id)
// 	}

// 	key0Succ, err = chord1.findSuccessor(ba0)
// 	if !bytes.Equal(key0Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key0Succ.Id)
// 	}

// 	key0Succ, err = chord3.findSuccessor(ba0)
// 	if !bytes.Equal(key0Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key0Succ.Id)
// 	}

// 	// FindSuccessor of key 1, call three servers, result should always be 1
// 	key1Succ, err := chord0.findSuccessor(ba1)
// 	if !bytes.Equal(key1Succ.Id, intToByteArray(1)) {
// 		t.Errorf("Find successor from node %d got = %d; want 1", chord0.GetId(), key1Succ.Id)
// 	}

// 	key1Succ, err = chord1.findSuccessor(ba1)
// 	if !bytes.Equal(key1Succ.Id, intToByteArray(1)) {
// 		t.Errorf("Find successor from node %d got = %d; want 1", chord1.GetId(), key1Succ.Id)
// 	}

// 	key1Succ, err = chord3.findSuccessor(ba1)
// 	if !bytes.Equal(key1Succ.Id, intToByteArray(1)) {
// 		t.Errorf("Find successor from node %d got = %d; want 1", chord3.GetId(), key1Succ.Id)
// 	}

// 	// FindSuccessor of key 2, call three servers, result should always be 3
// 	key2Succ, err := chord0.findSuccessor(ba2)
// 	if !bytes.Equal(key2Succ.Id, intToByteArray(2)) {
// 		t.Errorf("Find successor from node %d got = %d; want 2", chord0.GetId(), key2Succ.Id)
// 	}

// 	key2Succ, err = chord1.findSuccessor(ba2)
// 	if !bytes.Equal(key2Succ.Id, intToByteArray(2)) {
// 		t.Errorf("Find successor from node %d got = %d; want 2", chord1.GetId(), key2Succ.Id)
// 	}

// 	key2Succ, err = chord3.findSuccessor(ba2)
// 	if !bytes.Equal(key2Succ.Id, intToByteArray(2)) {
// 		t.Errorf("Find successor from node %d got = %d; want 2", chord3.GetId(), key2Succ.Id)
// 	}

// 	// FindSuccessor of key 3, call three servers, result should always be 3
// 	key3Succ, err := chord0.findSuccessor(ba3)
// 	if !bytes.Equal(key3Succ.Id, intToByteArray(3)) {
// 		t.Errorf("Find successor from node %d got = %d; want 3", chord0.GetId(), key3Succ.Id)
// 	}

// 	key2Succ, err = chord1.findSuccessor(ba3)
// 	if !bytes.Equal(key3Succ.Id, intToByteArray(3)) {
// 		t.Errorf("Find successor from node %d got = %d; want 3", chord1.GetId(), key3Succ.Id)
// 	}

// 	key2Succ, err = chord3.findSuccessor(ba3)
// 	if !bytes.Equal(key3Succ.Id, intToByteArray(3)) {
// 		t.Errorf("Find successor from node %d got = %d; want 3", chord3.GetId(), key3Succ.Id)
// 	}

// 	// FindSuccessor of key 4, call three servers, result should always be 0
// 	key4Succ, err := chord0.findSuccessor(ba4)
// 	if !bytes.Equal(key4Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key4Succ.Id)
// 	}

// 	key4Succ, err = chord1.findSuccessor(ba4)
// 	if !bytes.Equal(key4Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key4Succ.Id)
// 	}

// 	key4Succ, err = chord3.findSuccessor(ba4)
// 	if !bytes.Equal(key4Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key4Succ.Id)
// 	}

// 	// FindSuccessor of key 5, call three servers, result should always be 0
// 	key5Succ, err := chord0.findSuccessor(ba5)
// 	if !bytes.Equal(key5Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key5Succ.Id)
// 	}

// 	key5Succ, err = chord1.findSuccessor(ba5)
// 	if !bytes.Equal(key5Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key5Succ.Id)
// 	}

// 	key5Succ, err = chord3.findSuccessor(ba5)
// 	if !bytes.Equal(key5Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key5Succ.Id)
// 	}

// 	// FindSuccessor of key 6, call three servers, result should always be 0
// 	key6Succ, err := chord0.findSuccessor(ba6)
// 	if !bytes.Equal(key6Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key6Succ.Id)
// 	}

// 	key6Succ, err = chord1.findSuccessor(ba6)
// 	if !bytes.Equal(key6Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key6Succ.Id)
// 	}

// 	key6Succ, err = chord3.findSuccessor(ba6)
// 	if !bytes.Equal(key6Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key6Succ.Id)
// 	}

// 	// FindSuccessor of key 7, call three servers, result should always be 0
// 	key7Succ, err := chord0.findSuccessor(ba7)
// 	if !bytes.Equal(key7Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord0.GetId(), key7Succ.Id)
// 	}

// 	key7Succ, err = chord1.findSuccessor(ba7)
// 	if !bytes.Equal(key7Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord1.GetId(), key7Succ.Id)
// 	}

// 	key7Succ, err = chord3.findSuccessor(ba7)
// 	if !bytes.Equal(key7Succ.Id, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", chord3.GetId(), key6Succ.Id)
// 	}

// 	DPrintf("Server0:\n %s", chord0.tracer.String())
// 	DPrintf("Server1:\n %s", chord1.tracer.String())
// 	DPrintf("Server2:\n %s", server2.tracer.String())
// 	DPrintf("Server3:\n %s", chord3.tracer.String())
// }
