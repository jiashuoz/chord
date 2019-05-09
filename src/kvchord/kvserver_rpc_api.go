package kvchord

import (
	"context"
	"github.com/jiashuoz/chord/kvrpc"
	"google.golang.org/grpc"
	"log"
	"time"
)

// rpcConnWrapper contains grpc.ClientConn and chordrpc.ChordClient
type rpcConnWrapper struct {
	conn     *grpc.ClientConn
	kvClient kvrpc.KVClient
}

// Dial wraps around grpc Dial
func Dial(ip string) (*grpc.ClientConn, error) {
	return grpc.Dial(ip,
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.FailOnNonTempDialError(true),
		grpc.WithInsecure(),
	)
}

func (kv *KVServer) getRPCConn(ip string) kvrpc.KVClient {
	kv.rpcConnWrappersRWMu.RLock()
	rpcConn, ok := kv.rpcConnWrappers[ip]
	kv.rpcConnWrappersRWMu.RUnlock()
	if ok {
		return rpcConn.kvClient
	}
	return nil
}

func (kv *KVServer) connectRemote(ip string) (kvrpc.KVClient, error) {
	kvClient := kv.getRPCConn(ip)

	if kvClient != nil {
		return kvClient, nil
	}

	conn, err := Dial(ip)
	if err != nil {
		log.Fatal("connectRemote", err) // if any Dial error, crash the program
		return nil, err
	}

	kvClient = kvrpc.NewKVClient(conn)
	return kvClient, nil
}

// Get returns a value
func (kv *KVServer) getRPC(remoteIP string, key string) (string, error) {
	client, _ := kv.connectRemote(remoteIP)

	request := &kvrpc.GetRequest{Key: key}
	result, err := client.Get(context.Background(), request)
	return result.Val, err
}

func (kv *KVServer) putRPC(remoteIP string, key string, val string) error {
	client, _ := kv.connectRemote(remoteIP)

	request := &kvrpc.PutRequest{Key: key, Val: val}
	_, err := client.Put(context.Background(), request)
	return err
}

func (kv *KVServer) deleteRPC(remoteIP string, key string) (string, error) {
	client, _ := kv.connectRemote(remoteIP)

	request := &kvrpc.DeleteRequest{Key: key}
	result, err := client.Delete(context.Background(), request)
	return result.Val, err
}
