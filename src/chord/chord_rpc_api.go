package chord

import (
	"context"
	"github.com/jiashuoz/chord/chordrpc"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/benchmark"
	"time"
)

// A wrapper around some grpc internal connections
type grpcConn struct {
	ip     string
	client chordrpc.ChordClient // Chord service client
	conn   *grpc.ClientConn     // grpc client for underlying connection
}

// Dial returns a grpc.ClientConn
func Dial(ip string, options ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// return benchmark.NewClientConnWithContext(ctx, ip, options...), nil

	// return grpc.Dial(ip, options...)
	return grpc.DialContext(ctx,
		ip,
		options...,
	)
}

func (chord *ChordServer) connectRemote(remoteIP string) (chordrpc.ChordClient, error) {
	// time.Sleep(30 * time.Millisecond)
	chord.connectionsPoolRWMu.RLock()
	grpcc, ok := chord.connectionsPool[remoteIP]
	if ok {
		chord.connectionsPoolRWMu.RUnlock()
		return grpcc.client, nil
	}
	chord.connectionsPoolRWMu.RUnlock()

	conn, err := Dial(remoteIP, chord.config.DialOptions...)
	if err != nil {
		chord.logger.Println("connectRemote: Dial: ")
		return nil, err
	}

	client := chordrpc.NewChordClient(conn)
	grpcc = &grpcConn{remoteIP, client, conn}

	chord.connectionsPoolRWMu.Lock()
	chord.connectionsPool[remoteIP] = grpcc
	chord.connectionsPoolRWMu.Unlock()
	return client, nil
}

// findSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) notifyRPC(remote *chordrpc.Node, potentialPred *chordrpc.Node) (*chordrpc.NN, error) {
	client, err := chord.connectRemote(remote.Ip)
	if err != nil {
		return nil, err
	}

	result, err := client.Notify(context.Background(), potentialPred)
	return result, err
}

// findSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) findSuccessorRPC(remote *chordrpc.Node, id []byte) (*chordrpc.Node, error) {
	client, err := chord.connectRemote(remote.Ip)
	if err != nil {
		return nil, err
	}

	result, err := client.FindSuccessor(context.Background(), &chordrpc.ID{Id: id})
	return result, err
}

// findClosestPrecedingNodeRPC sends RPC call to remote node, returns closest node based on id
func (chord *ChordServer) findClosestPrecedingNodeRPC(remote *chordrpc.Node, id []byte) (*chordrpc.Node, error) {
	client, err := chord.connectRemote(remote.Ip)
	if err != nil {
		return nil, err
	}

	result, err := client.FindClosestPrecedingNode(context.Background(), &chordrpc.ID{Id: id})
	return result, err
}

// GetSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) getSuccessorRPC(remote *chordrpc.Node) (*chordrpc.Node, error) {
	client, err := chord.connectRemote(remote.Ip)
	if err != nil {
		return nil, err
	}

	result, err := client.GetSuccessor(context.Background(), &chordrpc.NN{})
	return result, err
}

// GetPredecessorRPC sends RPC call to remote node
func (chord *ChordServer) getPredecessorRPC(remote *chordrpc.Node) (*chordrpc.Node, error) {
	client, err := chord.connectRemote(remote.Ip)
	if err != nil {
		return nil, err
	}

	result, err := client.GetPredecessor(context.Background(), &chordrpc.NN{})
	return result, err
}

func (chord *ChordServer) setPredecessorRPC(remote *chordrpc.Node, pred *chordrpc.Node) (*chordrpc.NN, error) {
	client, err := chord.connectRemote(remote.Ip)
	if err != nil {
		return nil, err
	}

	result, err := client.SetPredecessor(context.Background(), pred)
	return result, err
}
func (chord *ChordServer) setSuccessorRPC(remote *chordrpc.Node, succ *chordrpc.Node) (*chordrpc.NN, error) {
	client, err := chord.connectRemote(remote.Ip)
	if err != nil {
		return nil, err
	}

	result, err := client.SetSuccessor(context.Background(), succ)
	return result, err
}

// Get returns a value
func (chord *ChordServer) getRPC(remoteIP string, key string) (string, error) {
	client, err := chord.connectRemote(remoteIP)
	if err != nil {
		return "", err
	}

	request := &chordrpc.GetRequest{Key: key}
	result, err := client.Get(context.Background(), request)
	return result.Val, err
}

func (chord *ChordServer) putRPC(remoteIP string, key string, val string) error {
	client, err := chord.connectRemote(remoteIP)
	if err != nil {
		return err
	}

	request := &chordrpc.PutRequest{Key: key, Val: val}
	_, err = client.Put(context.Background(), request)
	return err
}

func (chord *ChordServer) deleteRPC(remoteIP string, key string) (string, error) {
	client, err := chord.connectRemote(remoteIP)
	if err != nil {
		return "", err
	}

	request := &chordrpc.DeleteRequest{Key: key}
	result, err := client.Delete(context.Background(), request)
	return result.Val, err
}

func (chord *ChordServer) requestKeysRPC(succ *chordrpc.Node, start []byte, end []byte) ([]*chordrpc.KeyValue, error) {
	client, err := chord.connectRemote(succ.Ip)
	if err != nil {
		return nil, err
	}

	request := &chordrpc.RequestKeyValueRequest{Start: start, End: end}
	result, err := client.RequestKeys(context.Background(), request)
	return result.KeyValues, err
}

// func (chord *ChordServer) transferKeysRPC(succ *chordrpc.Node, start []byte, end []byte) error {
// 	client, err := chord.connectRemote(succ.Ip)
// 	if err != nil {
// 		return err
// 	}

// 	chord.kvStoreRWMu.RLock()
// 	kvs := make([]*chordrpc.KeyValue, 0)
// 	for k, v := range chord.kvStore.storage {
// 		hashedKey := chord.Hash(k)
// 		if betweenRightInclusive(hashedKey, start, end) {
// 			kvpair := &chordrpc.KeyValue{
// 				Key: k,
// 				Val: v,
// 			}
// 			kvs = append(kvs, kvpair)
// 		}
// 	}
// 	chord.kvStoreRWMu.RUnlock()
// 	request := &chordrpc.RequestKeyValueReply{KeyValues: kvs}
// 	result, err := client.TransferKeys(context.Background(), request)
// 	return err
// }

// func (chord *ChordServer) startCleanupConn() {
// 	ticker := time.NewTicker(5 * time.Second)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			chord.cleanupConnPool()
// 		}
// 	}
// }

// func (chord *ChordServer) cleanupConnPool() {
// 	chord.connectionsPoolRWMu.Lock()
// 	defer chord.connectionsPoolRWMu.Unlock()
// 	for host, grpcC := range chord.connectionsPool {
// 		if time.Since(grpcC.lastActive) > maxIdleTime {
// 			grpcC.conn.Close()
// 			delete(chord.connectionsPool, host)
// 		}
// 	}
// }
