package chord

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/jiashuoz/chord/chordrpc"
	"google.golang.org/grpc"
	"math/big"
	"math/rand"
	"net"
	"sync"
	"time"
)

const stabilizeTime = 50 * time.Millisecond
const fixFingerTime = 50 * time.Millisecond
const numBits = 3

type ChordServer struct {
	*chordrpc.Node // we never change this node, so no need to lock!!!

	predecessor *chordrpc.Node
	predRWMu    sync.RWMutex

	fingerTable []*chordrpc.Node
	fingerRWMu  sync.RWMutex

	grpcServer *grpc.Server

	stopChan chan struct{} // struct is the smallest thing in go, memory usage low

	rpcConnWrappers     map[string]*rpcConnWrapper // reuse existing conn to other servers
	rpcConnWrappersRWMu sync.RWMutex
}

func MakeChord(ip string, joinNode *chordrpc.Node) (*ChordServer, error) {
	chord := &ChordServer{
		Node: new(chordrpc.Node),
	}
	chord.Ip = ip
	chord.Id = hash(ip)
	chord.fingerTable = make([]*chordrpc.Node, numBits)
	chord.stopChan = make(chan struct{})
	chord.rpcConnWrappers = make(map[string]*rpcConnWrapper)

	listener, err := net.Listen("tcp", ip)
	if err != nil {
		checkError("", err)
		return nil, err
	}

	chord.grpcServer = grpc.NewServer()
	chordrpc.RegisterChordServer(chord.grpcServer, chord)
	// gmajpb.RegisterGMajServer(node.grpcs, node)

	go chord.grpcServer.Serve(listener)

	err = chord.join(joinNode)
	if err != nil {
		checkError("", err)
		return nil, err
	}

	go func() {
		ticker := time.NewTicker(stabilizeTime)
		for {
			select {
			case <-ticker.C:
				go chord.stabilize()
			case <-chord.stopChan:
				ticker.Stop()
				DPrintf("%v Stopping stabilize", chord.Id)
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(fixFingerTime)
		for {
			select {
			case <-ticker.C:
				go chord.fixFingers()
			case <-chord.stopChan:
				ticker.Stop()
				DPrintf("%v Stopping fixFingers", chord.Id)
				return
			}
		}
	}()

	return chord, nil
}

func (chord *ChordServer) Stop(params ...string) {
	fmt.Printf("%v is stopping....\n", chord.Id)
	for _, stopMsg := range params {
		fmt.Println(stopMsg)
	}
	// Stop all goroutines
	close(chord.stopChan)

	for _, rpcConnWrapper := range chord.rpcConnWrappers {
		err := rpcConnWrapper.conn.Close()
		checkError("Stop", err)
	}
}

func (chord *ChordServer) join(joinNode *chordrpc.Node) error {
	chord.predecessor = nil
	if joinNode == nil { // only node in the ring
		chord.fingerTable[0] = chord.Node
		// chord.predecessor = chord.Node
		return nil
	}
	succ, err := chord.findSuccessorRPC(joinNode, chord.Id)
	if idsEqual(succ.Id, chord.Id) {
		chord.Stop("Node with same ID already exists in the ring")
		return errors.New("Node with same ID already exists in the ring")
	}
	checkError("join", err)
	chord.fingerRWMu.Lock()
	chord.fingerTable[0] = succ
	chord.fingerRWMu.Unlock()
	return nil
}

// periodically verify's chord's immediate successor
func (chord *ChordServer) stabilize() error {
	succ := chord.getSuccessor()
	x, err := chord.getPredecessorRPC(succ)

	// RPC fails
	if err != nil {
		DPrintf("stabilize error: %v", err)
		return err
	}
	// the pred of our succ is nil, it hasn't updated it pred, still notify
	if x.Id == nil {
		chord.notifyRPC(chord.getSuccessor(), chord.Node)
		return nil
	}

	// found new successor
	if between(x.Id, chord.Id, succ.Id) {
		chord.fingerRWMu.Lock()
		chord.fingerTable[0] = x
		chord.fingerRWMu.Unlock()
	}

	chord.notifyRPC(chord.getSuccessor(), chord.Node)
	return nil
}

// notify tells chord, potentialPred thinks it might be chord's predecessor.
func (chord *ChordServer) notify(potentialPred *chordrpc.Node) error {
	chord.predRWMu.Lock()
	defer chord.predRWMu.Unlock()
	if chord.predecessor == nil || between(potentialPred.Id, chord.predecessor.Id, chord.Id) {
		chord.predecessor = potentialPred
	}
	return nil
}

// periodically refresh finger table entries
func (chord *ChordServer) fixFingers() {
	i := rand.Intn(numBits-1) + 1
	fingerStart := chord.fingerStart(i)
	finger, err := chord.findSuccessor(fingerStart)
	// Either RPC fails or something fails...
	if err != nil {
		DPrintf("fixFingers %v", err)
		return
	}
	chord.fingerRWMu.Lock()
	chord.fingerTable[i] = finger
	chord.fingerRWMu.Unlock()
}

// helper function used by fixFingers
func (chord *ChordServer) fingerStart(i int) []byte {
	currID := new(big.Int).SetBytes(chord.Id)
	offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
	maxVal := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(numBits), nil)
	start := new(big.Int).Add(currID, offset)
	start.Mod(start, maxVal)
	if len(start.Bytes()) == 0 {
		return []byte{0}
	}
	return start.Bytes()
}

func (chord *ChordServer) findSuccessor(id []byte) (*chordrpc.Node, error) {
	pred, err := chord.findPredecessor(id)
	if err != nil {
		return nil, err
	}

	succ, err := chord.getSuccessorRPC(pred)
	if err != nil {
		return nil, err
	}

	return succ, nil
}

func (chord *ChordServer) findPredecessor(id []byte) (*chordrpc.Node, error) {
	closest := chord.findClosestPrecedingNode(id)
	if idsEqual(closest.Id, chord.Id) {
		return closest, nil
	}

	// Get closest's successor
	closestSucc, _ := chord.getSuccessorRPC(closest)

	for !betweenRightInclusive(id, closest.Id, closestSucc.Id) {
		var err error
		closest, err = chord.findClosestPrecedingNodeRPC(closest, id)
		if err != nil {
			return nil, err
		}

		closestSucc, err = chord.getSuccessorRPC(closest) // get closest's successor
		if err != nil {
			return nil, err
		}
	}

	return closest, nil
}

func (chord *ChordServer) findClosestPrecedingNode(id []byte) *chordrpc.Node {
	chord.fingerRWMu.RLock()
	defer chord.fingerRWMu.RUnlock()
	for i := numBits - 1; i >= 0; i-- {
		if chord.fingerTable[i] != nil {
			if between(chord.fingerTable[i].Id, chord.Id, id) {
				return chord.fingerTable[i]
			}
		}
	}
	return chord.Node
}

func (chord *ChordServer) getSuccessor() *chordrpc.Node {
	chord.fingerRWMu.RLock()
	defer chord.fingerRWMu.RUnlock()
	return chord.fingerTable[0]
}

func (chord *ChordServer) getPredecessor() *chordrpc.Node {
	chord.predRWMu.RLock()
	defer chord.predRWMu.RUnlock()
	return chord.predecessor
}

func hash(ipAddr string) []byte {
	h := sha1.New()
	h.Write([]byte(ipAddr))

	idInt := big.NewInt(0)
	idInt.SetBytes(h.Sum(nil)) // Sum() returns []byte, convert it into BigInt

	maxVal := big.NewInt(0)
	maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil) // calculate 2^m
	idInt.Mod(idInt, maxVal)                            // mod id to make it to be [0, 2^m - 1]
	if idInt.Cmp(big.NewInt(0)) == 0 {
		return []byte{0}
	}
	return idInt.Bytes()
}

func (chord *ChordServer) String() string {
	chord.fingerRWMu.RLock()
	chord.predRWMu.RLock()
	defer chord.fingerRWMu.RUnlock()
	defer chord.predRWMu.RUnlock()
	str := chord.Node.String() + "\n"

	str += "Finger table:\n"
	str += "ith | start | successor\n"
	for i := 0; i < numBits; i++ {
		if chord.fingerTable[i] == nil {
			continue
		}
		currID := new(big.Int).SetBytes(chord.Id)
		offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
		maxVal := big.NewInt(0)
		maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil)
		start := new(big.Int).Add(currID, offset)
		start.Mod(start, maxVal)
		successor := chord.fingerTable[i].Id
		str += fmt.Sprintf("%d   | %d     | %d\n", i, start, successor)
	}

	str += fmt.Sprintf("Predecessor: ")
	if chord.predecessor == nil {
		str += fmt.Sprintf("none")
	} else {
		str += fmt.Sprintf("%v", chord.predecessor.Id)
	}
	return str
}
