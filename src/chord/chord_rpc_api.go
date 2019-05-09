package chord

import (
	"context"
	"github.com/jiashuoz/chord/chordrpc"
	"google.golang.org/grpc"
	"time"
)

// rpcConnWrapper contains grpc.ClientConn and chordrpc.ChordClient
type rpcConnWrapper struct {
	conn        *grpc.ClientConn
	chordClient chordrpc.ChordClient
}

func Dial(ip string) (*grpc.ClientConn, error) {
	return grpc.Dial(ip,
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
		grpc.FailOnNonTempDialError(true),
		grpc.WithInsecure(),
	)
}

func (chord *ChordServer) getRPCConn(ip string) chordrpc.ChordClient {
	chord.rpcConnWrappersRWMu.RLock()
	rpcConn, ok := chord.rpcConnWrappers[ip]
	chord.rpcConnWrappersRWMu.RUnlock()
	if ok {
		return rpcConn.chordClient
	}
	return nil
}

func (chord *ChordServer) connectRemote(remote *chordrpc.Node) (chordrpc.ChordClient, error) {
	ip := remote.Ip
	chordClient := chord.getRPCConn(ip)

	if chordClient != nil {
		return chordClient, nil
	}

	conn, err := Dial(ip)
	if err != nil {
		checkError("connectRemote", err) // if any Dial error, crash the program
		return nil, err
	}

	chordClient = chordrpc.NewChordClient(conn)

	newRPCConn := &rpcConnWrapper{conn, chordClient}
	chord.rpcConnWrappersRWMu.Lock()
	chord.rpcConnWrappers[ip] = newRPCConn
	chord.rpcConnWrappersRWMu.Unlock()
	return chordClient, nil
}

// findSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) notifyRPC(remote *chordrpc.Node, potentialPred *chordrpc.Node) (*chordrpc.NN, error) {
	client, _ := chord.connectRemote(remote)

	result, err := client.Notify(context.Background(), potentialPred)
	checkErrorGrace("Notify", err)
	return result, err
}

// findSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) findSuccessorRPC(remote *chordrpc.Node, id []byte) (*chordrpc.Node, error) {
	client, _ := chord.connectRemote(remote)

	result, err := client.FindSuccessor(context.Background(), &chordrpc.ID{Id: id})
	checkErrorGrace("FindSuccessor", err)
	return result, err
}

// findClosestPrecedingNodeRPC sends RPC call to remote node, returns closest node based on id
func (chord *ChordServer) findClosestPrecedingNodeRPC(remote *chordrpc.Node, id []byte) (*chordrpc.Node, error) {
	client, _ := chord.connectRemote(remote)

	result, err := client.FindClosestPrecedingNode(context.Background(), &chordrpc.ID{Id: id})
	checkErrorGrace("findClosestPrecedingNodeRPC", err)
	return result, err
}

// GetSuccessorRPC sends RPC call to remote node
func (chord *ChordServer) getSuccessorRPC(remote *chordrpc.Node) (*chordrpc.Node, error) {
	client, _ := chord.connectRemote(remote)

	result, err := client.GetSuccessor(context.Background(), &chordrpc.NN{})
	checkErrorGrace("getSuccessorRPC", err)
	return result, err
}

// GetPredecessorRPC sends RPC call to remote node
func (chord *ChordServer) getPredecessorRPC(remote *chordrpc.Node) (*chordrpc.Node, error) {
	client, _ := chord.connectRemote(remote)

	result, err := client.GetPredecessor(context.Background(), &chordrpc.NN{})
	checkErrorGrace("getPredecessorRPC", err)
	return result, err
}
