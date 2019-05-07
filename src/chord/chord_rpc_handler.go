package chord

import (
	"context"
	"github.com/jiashuoz/chord/chordrpc"
)

// Notify notifies Chord that potentialPred thinks it might be the predecessor
func (chord *ChordServer) Notify(ctx context.Context, potentialPred *chordrpc.Node) (*chordrpc.NN, error) {
	// need lock
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
	succ := chord.fingerTable[0]

	if succ == nil {
		return &chordrpc.Node{}, nil
	}

	return succ, nil
}

// GetPredecessor gets the predecessor of the node.
func (chord *ChordServer) GetPredecessor(context.Context, *chordrpc.NN) (*chordrpc.Node, error) {
	// need lock
	pred := chord.predecessor

	if pred == nil {
		return &chordrpc.Node{}, nil
	}

	return pred, nil
}

// SetPredecessor sets predecessor for chord
func (chord *ChordServer) SetPredecessor(ctx context.Context, pred *chordrpc.Node) (*chordrpc.NN, error) {
	// need lock
	chord.predecessor = pred

	return &chordrpc.NN{}, nil
}

// SetSuccessor sets predecessor for chord
func (chord *ChordServer) SetSuccessor(ctx context.Context, succ *chordrpc.Node) (*chordrpc.NN, error) {
	// need lock
	chord.fingerTable[0] = succ

	return &chordrpc.NN{}, nil
}
