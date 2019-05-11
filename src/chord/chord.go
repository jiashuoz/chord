package chord

import (
	"errors"
	"fmt"
	"github.com/jiashuoz/chord/chordrpc"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/benchmark"
	"google.golang.org/grpc/benchmark/latency"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type ChordServer struct {
	*chordrpc.Node // we never change this node, so no need to lock!!!

	predecessor     *chordrpc.Node
	predecessorRWMu sync.RWMutex

	fingerTable     []*chordrpc.Node
	fingerTableRWMu sync.RWMutex

	successorList     []*chordrpc.Node // keep track of log(n) nearest successor for recovery
	successorListRWMu sync.RWMutex

	stopChan chan struct{} // struct is the smallest thing in go, memory usage low

	kvStore     *kvStore
	kvStoreRWMu sync.RWMutex

	grpcServer          *grpc.Server         // this provides services
	connectionsPool     map[string]*grpcConn // this takes care of all connections
	connectionsPoolRWMu sync.RWMutex

	tracer     *Tracer // for testing latency and hops
	tracerRWMu sync.RWMutex

	config *Config

	logger *log.Logger

	network *latency.Network
}

// MakeChord takes an ip address, the ip of an existing node, returns a new ChordServer instance
func MakeChord(config *Config, ip string, joinNode string) (*ChordServer, error) {
	chord := &ChordServer{
		Node:   new(chordrpc.Node),
		config: config,
	}
	chord.Ip = ip
	chord.Id = chord.Hash(ip)
	chord.fingerTable = make([]*chordrpc.Node, chord.config.ringSize)
	chord.stopChan = make(chan struct{})
	chord.connectionsPool = make(map[string]*grpcConn)
	chord.kvStore = NewKVStore()

	chord.tracer = MakeTracer()

	chord.logger = log.New(os.Stderr, "logger: ", log.Ltime|log.Lshortfile) // print time and line number

	chord.network = &latency.Network{100 * 1024, 2 * time.Millisecond, 1500}
	l, err := net.Listen("tcp", ip)
	if err != nil {
		chord.logger.Println(err)
		return nil, err
	}
	// listener := chord.network.Listener(l) // listener with latency injected

	chord.grpcServer = grpc.NewServer()

	chordrpc.RegisterChordServer(chord.grpcServer, chord)

	// info := benchmark.ServerInfo{Type: "protobuf", Listener: l}

	// benchmark.StartServer(info)
	go chord.grpcServer.Serve(l)

	err = chord.join(&chordrpc.Node{Ip: joinNode})

	if err != nil {
		chord.logger.Println(err)
		return nil, err
	}

	go func() {
		ticker := time.NewTicker(chord.config.stabilizeTime)
		for {
			select {
			case <-ticker.C:
				err := chord.stabilize()
				if err != nil {
					chord.logger.Println(err)
				}
			case <-chord.stopChan:
				ticker.Stop()
				chord.logger.Printf("%s %d stopping stabilize\n", chord.Ip, chord.Id)
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(chord.config.fixFingerTime)
		i := 0
		for {
			select {
			case <-ticker.C:
				i, err = chord.fixFingers(i)
				if err != nil {
					chord.logger.Println(err)
				}
			case <-chord.stopChan:
				ticker.Stop()
				chord.logger.Printf("%s %d stopping fixing fingers\n", chord.Ip, chord.Id)
				return
			}
		}
	}()

	return chord, nil
}

// Lookup takes in a key, returns the ip address of the node that should store that chord pair
// return err is failure
func (chord *ChordServer) Lookup(key string) (string, error) {
	keyHased := chord.Hash(key)
	succ, err := chord.findSuccessor(keyHased) // check where this key should go in the ring
	if err != nil {                            // if failure, notify the application level
		chord.logger.Println(err)
		return "", err
	}
	return succ.Ip, nil
}

// Stop gracefully stops the chord instance
func (chord *ChordServer) Stop() {
	chord.logger.Printf("%s stopping instance\n", chord.Ip)
	// Stop all goroutines
	close(chord.stopChan)
	chord.connectionsPoolRWMu.Lock()
	defer chord.connectionsPoolRWMu.Unlock()
	for _, connection := range chord.connectionsPool {
		err := connection.conn.Close()
		if err != nil {
			chord.logger.Println(err)
		}
	}
}

func (chord *ChordServer) StopFixFingers() {
	close(chord.stopChan)
}

func (chord *ChordServer) join(joinNode *chordrpc.Node) error {
	chord.predecessor = nil
	if joinNode.Ip == "" { // only node in the ring
		chord.fingerTable[0] = chord.Node
		return nil
	}
	succ, err := chord.findSuccessorRPC(joinNode, chord.Id)
	if err != nil {
		return err // this will be logged by logger in MakeChord
	}
	if idsEqual(succ.Id, chord.Id) {
		return errors.New("Node with same ID already exists in the ring") // this will be logged by logger in MakeChord
	}
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
		chord.logger.Println("RPC fails or the successor node died...")
		return err
	}

	// the pred of our succ is nil, it hasn't updated it pred, still notify
	if x.Id == nil {
		_, err = chord.notifyRPC(chord.getSuccessor(), chord.Node)
		return err
	}

	// found new successor
	if between(x.Id, chord.Id, succ.Id) {
		chord.fingerTableRWMu.Lock()
		chord.fingerTable[0] = x
		chord.fingerTableRWMu.Unlock()
	}

	_, err = chord.notifyRPC(chord.getSuccessor(), chord.Node)
	return err
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
func (chord *ChordServer) fixFingers(i int) (int, error) {
	i = (i + 1) % chord.config.ringSize
	fingerStart := chord.fingerStart(i)
	finger, err := chord.findSuccessorForFingers(fingerStart)
	// Either RPC fails or something fails...
	if err != nil {
		return 0, err
	}
	chord.fingerTableRWMu.Lock()
	chord.fingerTable[i] = finger
	chord.fingerTableRWMu.Unlock()
	return i, nil
}

// helper function used by fixFingers
func (chord *ChordServer) fingerStart(i int) []byte {
	currID := new(big.Int).SetBytes(chord.Id)
	offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
	maxVal := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(chord.config.ringSize)), nil)
	start := new(big.Int).Add(currID, offset)
	start.Mod(start, maxVal)
	if len(start.Bytes()) == 0 {
		return []byte{0}
	}
	return start.Bytes()
}

// return the successor of id
func (chord *ChordServer) findSuccessorForFingers(id []byte) (*chordrpc.Node, error) {
	pred, err := chord.findPredecessorForFingers(id)
	if err != nil {
		chord.logger.Println("findSuccessor")
		return nil, err
	}

	succ, err := chord.getSuccessorRPC(pred)
	if err != nil {
		chord.logger.Println("findSuccessor: getSuccessorRPC")
		return nil, err
	}
	return succ, nil
}

// return the predecessor of id, this could launch RPC calls to other node
func (chord *ChordServer) findPredecessorForFingers(id []byte) (*chordrpc.Node, error) {
	closest := chord.findClosestPrecedingNode(id)
	if idsEqual(closest.Id, chord.Id) {
		return closest, nil
	}

	// Get closest's successor
	closestSucc, err := chord.getSuccessorRPC(closest)
	if err != nil {
		chord.logger.Println("findPredecessor: getSuccessorRPC")
		return nil, err
	}

	for !betweenRightInclusive(id, closest.Id, closestSucc.Id) {
		var err error
		closest, err = chord.findClosestPrecedingNodeRPC(closest, id)
		if err != nil {
			chord.logger.Println("findPredecessor: findClosestPrecedingNodeRPC")
			return nil, err
		}

		closestSucc, err = chord.getSuccessorRPC(closest) // get closest's successor
		if err != nil {
			chord.logger.Println("findPredecessor: getSuccessorRPC")
			return nil, err
		}
	}

	return closest, nil
}

// return the successor of id
func (chord *ChordServer) findSuccessor(id []byte) (*chordrpc.Node, error) {
	chord.tracerRWMu.Lock()
	chord.tracer.startTracer(chord.Id, id)
	pred, err := chord.findPredecessor(id)
	if err != nil {
		chord.logger.Println("findSuccessor")
		return nil, err
	}

	succ, err := chord.getSuccessorRPC(pred)
	if err != nil {
		chord.logger.Println("findSuccessor: getSuccessorRPC")
		return nil, err
	}
	chord.tracer.traceNode(succ.Id)
	chord.tracer.endTracer(succ.Id)
	chord.tracerRWMu.Unlock()
	return succ, nil
}

// return the predecessor of id, this could launch RPC calls to other node
func (chord *ChordServer) findPredecessor(id []byte) (*chordrpc.Node, error) {
	closest := chord.findClosestPrecedingNode(id)
	if idsEqual(closest.Id, chord.Id) {
		return closest, nil
	}

	// Get closest's successor
	closestSucc, err := chord.getSuccessorRPC(closest)
	if err != nil {
		chord.logger.Println("findPredecessor: getSuccessorRPC")
		return nil, err
	}

	chord.tracer.traceNode(closest.Id)

	for !betweenRightInclusive(id, closest.Id, closestSucc.Id) {
		var err error
		closest, err = chord.findClosestPrecedingNodeRPC(closest, id)
		if err != nil {
			chord.logger.Println("findPredecessor: findClosestPrecedingNodeRPC")
			return nil, err
		}

		closestSucc, err = chord.getSuccessorRPC(closest) // get closest's successor
		if err != nil {
			chord.logger.Println("findPredecessor: getSuccessorRPC")
			return nil, err
		}

		chord.tracer.traceNode(closest.Id)
	}

	return closest, nil
}

// returns the closest finger based on fingerTable, doesn't do RPC call
func (chord *ChordServer) findClosestPrecedingNode(id []byte) *chordrpc.Node {
	chord.fingerTableRWMu.RLock()
	defer chord.fingerTableRWMu.RUnlock()
	for i := chord.config.ringSize - 1; i >= 0; i-- {
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
	for i := 0; i < chord.config.ringSize; i++ {
		if chord.fingerTable[i] == nil {
			continue
		}
		currID := new(big.Int).SetBytes(chord.Id)
		offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
		maxVal := big.NewInt(0)
		maxVal.Exp(big.NewInt(2), big.NewInt(int64(chord.config.ringSize)), nil)
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

// KVStore
// ************************************************************************
// ************************************************************************

func (chord *ChordServer) get(key string) (string, error) {
	ip, err := chord.Lookup(key) // ip of the node to store key
	checkError("get", err)
	if strings.Compare(ip, chord.Ip) == 0 {
		return chord.getVal(key)
	}
	val, err := chord.getRPC(ip, key)
	checkError("get", err)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (chord *ChordServer) put(key string, val string) error {
	ip, err := chord.Lookup(key) // ip of the node to store key
	checkError("put", err)
	if strings.Compare(ip, chord.Ip) == 0 {
		chord.putVal(key, val)
		return nil
	}
	err = chord.putRPC(ip, key, val)
	checkError("put", err)
	if err != nil {
		return err
	}
	return nil
}

func (chord *ChordServer) delete(key string) (string, error) {
	ip, err := chord.Lookup(key)
	checkError("delete", err)
	if strings.Compare(ip, chord.Ip) == 0 {
		val := chord.deleteVal(key)
		return val, nil
	}
	val, err := chord.deleteRPC(ip, key)
	checkError("delete", err)
	if err != nil {
		return "", err
	}
	return val, nil
}

func (chord *ChordServer) getVal(key string) (string, error) {
	chord.kvStoreRWMu.RLock()
	defer chord.kvStoreRWMu.RLock()
	val, ok := chord.kvStore.Get(key)
	if ok != nil {
		return "", ok
	}
	return val, nil
}

func (chord *ChordServer) putVal(key string, val string) {
	chord.kvStoreRWMu.Lock()
	defer chord.kvStoreRWMu.Unlock()
	chord.kvStore.Put(key, val)
}

func (chord *ChordServer) deleteVal(key string) string {
	chord.kvStoreRWMu.Lock()
	defer chord.kvStoreRWMu.Unlock()
	val, ok := chord.kvStore.Get(key)
	if ok != nil {
		return ""
	}
	chord.kvStore.Delete(key)
	return val
}

func (chord *ChordServer) keyCount() int {
	chord.kvStoreRWMu.RLock()
	defer chord.kvStoreRWMu.RUnlock()
	return chord.kvStore.KeyCount()
}

// Hash takes a string, returns hashed values in []byte
func (chord *ChordServer) Hash(ipAddr string) []byte {
	h := chord.config.Hash()
	h.Write([]byte(ipAddr))

	idInt := big.NewInt(0)
	idInt.SetBytes(h.Sum(nil)) // Sum() returns []byte, convert it into BigInt

	maxVal := big.NewInt(0)
	maxVal.Exp(big.NewInt(2), big.NewInt(int64(chord.config.ringSize)), nil) // calculate 2^m
	idInt.Mod(idInt, maxVal)                                                 // mod id to make it to be [0, 2^m - 1]
	if idInt.Cmp(big.NewInt(0)) == 0 {
		return []byte{0}
	}
	return idInt.Bytes()
}
