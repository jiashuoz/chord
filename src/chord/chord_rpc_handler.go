package chord

import (
	"context"
	"github.com/jiashuoz/chord/chordrpc"
)

// Notify notifies Chord that potentialPred thinks it might be the predecessor
func (chord *ChordServer) Notify(ctx context.Context, potentialPred *chordrpc.Node) (*chordrpc.NN, error) {
	chord.notify(potentialPred)

	return &chordrpc.NN{}, nil
}

// FindSuccessor handles remote request, finds the successor of id
func (chord *ChordServer) FindSuccessor(ctx context.Context, id *chordrpc.ID) (*chordrpc.Node, error) {
	succ, err := chord.findSuccessor(id.Id)
	if err != nil {
		return nil, err
	}
	return succ, nil
}

// FindClosestPrecedingNode handles remote request, finds the closest finger with given id
func (chord *ChordServer) FindClosestPrecedingNode(ctx context.Context, id *chordrpc.ID) (*chordrpc.Node, error) {
	closestPrecedingNode := chord.findClosestPrecedingNode(id.Id)
	return closestPrecedingNode, nil
}

// GetSuccessor gets the successor of the node.
func (chord *ChordServer) GetSuccessor(context.Context, *chordrpc.NN) (*chordrpc.Node, error) {
	// need lock
	succ := chord.getSuccessor()
	if succ == nil {
		return &chordrpc.Node{}, nil
	}

	return succ, nil
}

// GetPredecessor gets the predecessor of the node.
func (chord *ChordServer) GetPredecessor(context.Context, *chordrpc.NN) (*chordrpc.Node, error) {
	// need lock
	pred := chord.getPredecessor()

	// if pred == nil, simply return an empty node to indicate that, err should be nil, otherwise rpc won't work
	if pred == nil {
		return &chordrpc.Node{}, nil
	}

	return pred, nil
}

// SetPredecessor sets predecessor for chord
func (chord *ChordServer) SetPredecessor(ctx context.Context, pred *chordrpc.Node) (*chordrpc.NN, error) {
	// need lock
	chord.predecessorRWMu.Lock()
	defer chord.predecessorRWMu.Unlock()
	chord.predecessor = pred

	return &chordrpc.NN{}, nil
}

// SetSuccessor sets predecessor for chord
func (chord *ChordServer) SetSuccessor(ctx context.Context, succ *chordrpc.Node) (*chordrpc.NN, error) {
	// need lock
	chord.fingerTableRWMu.Lock()
	defer chord.fingerTableRWMu.Unlock()
	chord.fingerTable[0] = succ

	return &chordrpc.NN{}, nil
}

func (chord *ChordServer) Get(ctx context.Context, args *chordrpc.GetRequest) (*chordrpc.GetReply, error) {
	val, err := chord.getVal(args.Key)
	return &chordrpc.GetReply{Val: val}, err
}

func (chord *ChordServer) Put(ctx context.Context, args *chordrpc.PutRequest) (*chordrpc.PutReply, error) {
	chord.putVal(args.Key, args.Val)
	return &chordrpc.PutReply{}, nil
}

func (chord *ChordServer) Delete(ctx context.Context, args *chordrpc.DeleteRequest) (*chordrpc.DeleteReply, error) {
	val := chord.deleteVal(args.Key)
	return &chordrpc.DeleteReply{Val: val}, nil
}

// RequestKeys sends back the requested keys
func (chord *ChordServer) RequestKeys(ctx context.Context, args *chordrpc.RequestKeyValueRequest) (*chordrpc.RequestKeyValueReply, error) {
	chord.kvStoreRWMu.Lock()
	defer chord.kvStoreRWMu.Unlock()
	start := args.Start
	end := args.End
	kvs := make([]*chordrpc.KeyValue, 0)
	for k, v := range chord.kvStore.storage {
		hashedKey := chord.Hash(k)
		if betweenRightInclusive(hashedKey, start, end) {
			kvpair := &chordrpc.KeyValue{
				Key: k,
				Val: v,
			}
			kvs = append(kvs, kvpair)
		}
	}
	return &chordrpc.RequestKeyValueReply{KeyValues: kvs}, nil
}
