package kvchord

import (
	"chord"
	"github.com/jiashuoz/chord/kvrpc"
	"google.golang.org/grpc"
	"log"
	"strings"
	"sync"
)

// KVServer is a key value server running on top of chord
type KVServer struct {
	chord *chord.ChordServer

	storage     map[string]string
	storageRWMu sync.RWMutex

	grpcServer *grpc.Server
}

func (kv *KVServer) getVal(key string) string {
	kv.storageRWMu.RLock()
	defer kv.storageRWMu.RLock()
	val, ok := kv.storage[key]
	if !ok {
		return ""
	}
	return val
}

func (kv *KVServer) putVal(key string, val string) {
	kv.storageRWMu.Lock()
	defer kv.storageRWMu.Unlock()
	kv.storage[key] = val
}

func (kv *KVServer) deleteVal(key string) string {
	kv.storageRWMu.Lock()
	defer kv.storageRWMu.RUnlock()
	val, ok := kv.storage[key]
	if !ok {
		return ""
	}
	delete(kv.storage, key)
	return val
}

func (kv *KVServer) get(key string) (string, error) {
	ip, err := kv.chord.Lookup(key) // ip of the node to store key
	checkError("", err)
	if strings.Compare(ip, kv.chord.Ip) == 0 {
		return kv.storage[key], nil
	}
	val, err := kv.getRPC(ip, key)
	checkError("", err)
	return val, nil
}

func (kv *KVServer) put(key string, val string) error {
	ip, err := kv.chord.Lookup(key) // ip of the node to store key
	checkError("", err)
	if strings.Compare(ip, kv.chord.Ip) == 0 {
		kv.storage[key] = val
		return nil
	}
	err = kv.putRPC(ip, key, val)
	checkError("", err)
	return nil
}

func (kv *KVServer) delete(key string) (string, error) {
	return "", nil
}

// Kill should stop a server
func (kv *KVServer) Kill() {
}

// StartKVServer creates a new KVServer on top of Chord
func StartKVServer(ip string, joinNodeIP string) *KVServer {

	kv := new(KVServer)
	kv.storage = make(map[string]string)
	var err error
	kv.chord, err = chord.MakeChord(ip, chord.MakeJoinNode(joinNodeIP))

	kv.grpcServer = grpc.NewServer()
	kvrpc.RegisterKVServer(kv.grpcServer, kv)

	if err != nil {
		log.Fatal(err)
	}
	return kv
}
