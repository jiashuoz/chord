package kvchord

import (
	"context"
	"github.com/jiashuoz/chord/kvrpc"
)

// Get returns a value from kv store
func (kv *KVServer) Get(ctx context.Context, getRequest *kvrpc.GetRequest) (*kvrpc.GetReply, error) {
	val := kv.getVal(getRequest.Key)
	return &kvrpc.GetReply{Val: val}, nil
}

// Put puts a value into kv store
func (kv *KVServer) Put(ctx context.Context, putRequest *kvrpc.PutRequest) (*kvrpc.PutReply, error) {
	kv.putVal(putRequest.Key, putRequest.Val)
	return &kvrpc.PutReply{}, nil
}

// Delete deletes a val from kv store
func (kv *KVServer) Delete(ctx context.Context, deleteRequest *kvrpc.DeleteRequest) (*kvrpc.DeleteReply, error) {
	val := kv.deleteVal(deleteRequest.Key)
	return &kvrpc.DeleteReply{Val: val}, nil
}
