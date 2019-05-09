package chord

import (
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
const numBits = 160
const numBits4 = 4
const r = 4

type ChordServer struct {
	*chordrpc.Node // we never change this node, so no need to lock!!!

	predecessor     *chordrpc.Node
	predecessorRWMu sync.RWMutex

	fingerTable     []*chordrpc.Node
	fingerTableRWMu sync.RWMutex

	successorList     []*chordrpc.Node // keep track of log(n) nearest successor for recovery
	successorListRWMu sync.RWMutex

	grpcServer *grpc.Server

	stopChan chan struct{} // struct is the smallest thing in go, memory usage low

	rpcConnWrappers     map[string]*rpcConnWrapper // reuse existing conn to other servers
	rpcConnWrappersRWMu sync.RWMutex
}

func MakeChord(ip string, joinNode *chordrpc.Node, grpcs *grpc.Server) (*ChordServer, error) {
	chord := &ChordServer{
		Node: new(chordrpc.Node),
	}
	chord.Ip = ip
	chord.Id = Hash(ip)
	chord.fingerTable = make([]*chordrpc.Node, numBits)
	chord.stopChan = make(chan struct{})
	chord.rpcConnWrappers = make(map[string]*rpcConnWrapper)
	chord.successorList = make([]*chordrpc.Node, numBits)

	listener, err := net.Listen("tcp", ip)
	if err != nil {
		checkError("", err)
		return nil, err
	}

	chord.grpcServer = grpcs
	chordrpc.RegisterChordServer(chord.grpcServer, chord)

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

// Lookup takes in a key, returns the ip address of the node that should store that kv pair
// return err is failure
func (chord *ChordServer) Lookup(key string) (string, error) {
	keyHased := Hash(key)
	succ, err := chord.findSuccessor(keyHased) // check where this key should go in the ring
	if err != nil {                            // if failure, notify the application level
		return "", err
	}
	return succ.Ip, nil
}

// Stop gracefully stops the chord instance
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
		chord.successorList[0] = chord.Node
		// chord.predecessor = chord.Node
		return nil
	}
	succ, err := chord.findSuccessorRPC(joinNode, chord.Id)
	if idsEqual(succ.Id, chord.Id) {
		chord.Stop("Node with same ID already exists in the ring")
		return errors.New("Node with same ID already exists in the ring")
	}
	checkError("join", err)
	chord.fingerTableRWMu.Lock()
	chord.fingerTable[0] = succ
	chord.fingerTableRWMu.Unlock()
	return nil
}

// periodically verify's chord's immediate successor
func (chord *ChordServer) stabilize() error {
	succ := chord.getSuccessor()
	x, err := chord.getPredecessorRPC(succ)

	// RPC fails or the successor node died...
	if x == nil || err != nil {
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
		chord.fingerTableRWMu.Lock()
		chord.fingerTable[0] = x
		chord.fingerTableRWMu.Unlock()
	}

	chord.notifyRPC(chord.getSuccessor(), chord.Node)
	return nil
}

// notify tells chord, potentialPred thinks it might be chord's predecessor.
func (chord *ChordServer) notify(potentialPred *chordrpc.Node) error {
	chord.predecessorRWMu.Lock()
	defer chord.predecessorRWMu.Unlock()
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
	chord.fingerTableRWMu.Lock()
	chord.fingerTable[i] = finger
	chord.fingerTableRWMu.Unlock()
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

// return the successor of id
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

// return the predecessor of id, this could launch RPC calls to other node
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

// returns the closest finger based on fingerTable, doesn't do RPC call
func (chord *ChordServer) findClosestPrecedingNode(id []byte) *chordrpc.Node {
	chord.fingerTableRWMu.RLock()
	defer chord.fingerTableRWMu.RUnlock()
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
	chord.fingerTableRWMu.RLock()
	defer chord.fingerTableRWMu.RUnlock()
	return chord.fingerTable[0]
}

func (chord *ChordServer) getPredecessor() *chordrpc.Node {
	chord.predecessorRWMu.RLock()
	defer chord.predecessorRWMu.RUnlock()
	return chord.predecessor
}

// Returns the state of chord server, uses lock, for debugging
func (chord *ChordServer) String() string {
	chord.fingerTableRWMu.RLock()
	chord.predecessorRWMu.RLock()
	defer chord.fingerTableRWMu.RUnlock()
	defer chord.predecessorRWMu.RUnlock()
	str := fmt.Sprintf("id: %v ip: %v\n", chord.Id, chord.Ip)

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

// MakeJoinNode is mostly used for application level, takes the ip of JoinNode and make a node around it
func MakeJoinNode(ip string) *chordrpc.Node {
	if ip == "" {
		return nil
	}

	return &chordrpc.Node{Id: Hash(ip), Ip: ip}
}
