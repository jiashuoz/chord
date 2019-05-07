package chord

import (
	"fmt"
	"os"
	"testing"
)

func TestMake(t *testing.T) {
	fmt.Println("Test MakeChord")
	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
	chord0, err := MakeChord(testAddrs[0], nil)
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

// func TestGetSuccessor(t *testing.T) {
// 	testAddrs := reverseHash(numBits, "127.0.0.1", 5000)
// 	server0 := MakeServer(testAddrs[0]) // id 0
// 	server1 := MakeServer(testAddrs[1]) // id 1
// 	server2 := MakeServer(testAddrs[2]) // id 2
// 	server3 := MakeServer(testAddrs[3]) // id 3

// 	node0 := MakeNode(testAddrs[0])
// 	node1 := MakeNode(testAddrs[1])
// 	node2 := MakeNode(testAddrs[2])
// 	node3 := MakeNode(testAddrs[3])

// 	server0.fingerTable[0] = node1
// 	server0.fingerTable[1] = node2
// 	server0.fingerTable[2] = node0

// 	server1.fingerTable[0] = node2
// 	server1.fingerTable[1] = node3
// 	server1.fingerTable[2] = node0

// 	server2.fingerTable[0] = node3
// 	server2.fingerTable[1] = node0
// 	server2.fingerTable[2] = node0

// 	server3.fingerTable[0] = node0
// 	server3.fingerTable[1] = node0
// 	server3.fingerTable[2] = node0

// 	// run(server0)
// 	// run(server2)
// 	// run(server1)
// 	// run(server3)

// 	fmt.Println(server0.String(true))
// 	fmt.Println(server1.String(true))
// 	fmt.Println(server2.String(true))
// 	fmt.Println(server3.String(true))

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
// 	key0Succ := server0.FindSuccessor(ba0)
// 	if !bytes.Equal(key0Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server0.GetID(), key0Succ.ID)
// 	}

// 	key0Succ = server1.FindSuccessor(ba0)
// 	if !bytes.Equal(key0Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server1.GetID(), key0Succ.ID)
// 	}

// 	key0Succ = server3.FindSuccessor(ba0)
// 	if !bytes.Equal(key0Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server3.GetID(), key0Succ.ID)
// 	}

// 	// FindSuccessor of key 1, call three servers, result should always be 1
// 	key1Succ := server0.FindSuccessor(ba1)
// 	if !bytes.Equal(key1Succ.ID, intToByteArray(1)) {
// 		t.Errorf("Find successor from node %d got = %d; want 1", server0.GetID(), key1Succ.ID)
// 	}

// 	key1Succ = server1.FindSuccessor(ba1)
// 	if !bytes.Equal(key1Succ.ID, intToByteArray(1)) {
// 		t.Errorf("Find successor from node %d got = %d; want 1", server1.GetID(), key1Succ.ID)
// 	}

// 	key1Succ = server3.FindSuccessor(ba1)
// 	if !bytes.Equal(key1Succ.ID, intToByteArray(1)) {
// 		t.Errorf("Find successor from node %d got = %d; want 1", server3.GetID(), key1Succ.ID)
// 	}

// 	// FindSuccessor of key 2, call three servers, result should always be 3
// 	key2Succ := server0.FindSuccessor(ba2)
// 	if !bytes.Equal(key2Succ.ID, intToByteArray(2)) {
// 		t.Errorf("Find successor from node %d got = %d; want 2", server0.GetID(), key2Succ.ID)
// 	}

// 	key2Succ = server1.FindSuccessor(ba2)
// 	if !bytes.Equal(key2Succ.ID, intToByteArray(2)) {
// 		t.Errorf("Find successor from node %d got = %d; want 2", server1.GetID(), key2Succ.ID)
// 	}

// 	key2Succ = server3.FindSuccessor(ba2)
// 	if !bytes.Equal(key2Succ.ID, intToByteArray(2)) {
// 		t.Errorf("Find successor from node %d got = %d; want 2", server3.GetID(), key2Succ.ID)
// 	}

// 	// FindSuccessor of key 3, call three servers, result should always be 3
// 	key3Succ := server0.FindSuccessor(ba3)
// 	if !bytes.Equal(key3Succ.ID, intToByteArray(3)) {
// 		t.Errorf("Find successor from node %d got = %d; want 3", server0.GetID(), key3Succ.ID)
// 	}

// 	key2Succ = server1.FindSuccessor(ba3)
// 	if !bytes.Equal(key3Succ.ID, intToByteArray(3)) {
// 		t.Errorf("Find successor from node %d got = %d; want 3", server1.GetID(), key3Succ.ID)
// 	}

// 	key2Succ = server3.FindSuccessor(ba3)
// 	if !bytes.Equal(key3Succ.ID, intToByteArray(3)) {
// 		t.Errorf("Find successor from node %d got = %d; want 3", server3.GetID(), key3Succ.ID)
// 	}

// 	// FindSuccessor of key 4, call three servers, result should always be 0
// 	key4Succ := server0.FindSuccessor(ba4)
// 	if !bytes.Equal(key4Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server0.GetID(), key4Succ.ID)
// 	}

// 	key4Succ = server1.FindSuccessor(ba4)
// 	if !bytes.Equal(key4Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server1.GetID(), key4Succ.ID)
// 	}

// 	key4Succ = server3.FindSuccessor(ba4)
// 	if !bytes.Equal(key4Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server3.GetID(), key4Succ.ID)
// 	}

// 	// FindSuccessor of key 5, call three servers, result should always be 0
// 	key5Succ := server0.FindSuccessor(ba5)
// 	if !bytes.Equal(key5Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server0.GetID(), key5Succ.ID)
// 	}

// 	key5Succ = server1.FindSuccessor(ba5)
// 	if !bytes.Equal(key5Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server1.GetID(), key5Succ.ID)
// 	}

// 	key5Succ = server3.FindSuccessor(ba5)
// 	if !bytes.Equal(key5Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server3.GetID(), key5Succ.ID)
// 	}

// 	// FindSuccessor of key 6, call three servers, result should always be 0
// 	key6Succ := server0.FindSuccessor(ba6)
// 	if !bytes.Equal(key6Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server0.GetID(), key6Succ.ID)
// 	}

// 	key6Succ = server1.FindSuccessor(ba6)
// 	if !bytes.Equal(key6Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server1.GetID(), key6Succ.ID)
// 	}

// 	key6Succ = server3.FindSuccessor(ba6)
// 	if !bytes.Equal(key6Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server3.GetID(), key6Succ.ID)
// 	}

// 	// FindSuccessor of key 7, call three servers, result should always be 0
// 	key7Succ := server0.FindSuccessor(ba7)
// 	if !bytes.Equal(key7Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server0.GetID(), key7Succ.ID)
// 	}

// 	key7Succ = server1.FindSuccessor(ba7)
// 	if !bytes.Equal(key7Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server1.GetID(), key7Succ.ID)
// 	}

// 	key7Succ = server3.FindSuccessor(ba7)
// 	if !bytes.Equal(key7Succ.ID, intToByteArray(0)) {
// 		t.Errorf("Find successor from node %d got = %d; want 0", server3.GetID(), key6Succ.ID)
// 	}

// 	DPrintf("Server0:\n %s", server0.tracer.String())
// 	DPrintf("Server1:\n %s", server1.tracer.String())
// 	DPrintf("Server2:\n %s", server2.tracer.String())
// 	DPrintf("Server3:\n %s", server3.tracer.String())
// }
