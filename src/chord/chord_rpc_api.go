package chord

import (
	"context"
	"github.com/jiashuoz/chord/chordrpc"
	"google.golang.org/grpc"
	"time"
)

func Dial(ip string) (*grpc.ClientConn, error) {
	return grpc.Dial(ip,
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.FailOnNonTempDialError(true),
		grpc.WithInsecure(),
	)
}

func (chord *ChordServer) connectRemote(remote *chordrpc.Node) (chordrpc.ChordClient, error) {
	ip := remote.Ip

	conn, err := Dial(ip)
	if err != nil {
		return nil, err
	}

	chordClient := chordrpc.NewChordClient(conn)

	return chordClient, nil
}

// findSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) notifyRPC(remote *chordrpc.Node, potentialPred *chordrpc.Node) (*chordrpc.NN, error) {
	client, err := chord.connectRemote(remote)
	if err != nil {
		return nil, err
	}

	return client.Notify(context.Background(), potentialPred)
}

// findSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) findSuccessorRPC(remote *chordrpc.Node, id []byte) (*chordrpc.Node, error) {
	client, err := chord.connectRemote(remote)
	if err != nil {
		return nil, err
	}

	return client.FindSuccessor(context.Background(), &chordrpc.ID{Id: id})
}

// findClosestPrecedingNodeRPC sends RPC call to remote node, returns closest node based on id
func (chord *ChordServer) findClosestPrecedingNodeRPC(remote *chordrpc.Node, id []byte) (*chordrpc.Node, error) {
	client, err := chord.connectRemote(remote)
	if err != nil {
		return nil, err
	}

	return client.FindClosestPrecedingNode(context.Background(), &chordrpc.ID{Id: id})
}

// GetSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) getSuccessorRPC(remote *chordrpc.Node) (*chordrpc.Node, error) {
	client, err := chord.connectRemote(remote)
	if err != nil {
		return nil, err
	}

	return client.GetSuccessor(context.Background(), &chordrpc.NN{})
}

// GetPredecessorRPC sends RPC call to remote node
func (chord *ChordServer) getPredecessorRPC(remote *chordrpc.Node) (*chordrpc.Node, error) {
	client, err := chord.connectRemote(remote)
	if err != nil {
		return nil, err
	}

	return client.GetPredecessor(context.Background(), &chordrpc.NN{})
}
