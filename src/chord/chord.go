package chord

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/jiashuoz/chord/chordrpc"
	"google.golang.org/grpc"
	"math/big"
	"net"
)

const numBits = 3

type ChordServer struct {
	*chordrpc.Node

	predecessor *chordrpc.Node

	fingerTable []*chordrpc.Node

	grpcServer *grpc.Server
}

func MakeChord(ip string, joinNode *chordrpc.Node) (*ChordServer, error) {
	chord := &ChordServer{
		Node: new(chordrpc.Node),
	}
	chord.Ip = ip
	chord.Id = hash(ip)
	chord.fingerTable = make([]*chordrpc.Node, numBits)

	listener, err := net.Listen("tcp", ip)
	if err != nil {
		return nil, err
	}

	chord.grpcServer = grpc.NewServer()
	chordrpc.RegisterChordServer(chord.grpcServer, chord)
	// gmajpb.RegisterGMajServer(node.grpcs, node)

	go chord.grpcServer.Serve(listener)

	err = chord.join(joinNode)
	if err != nil {
		return nil, err
	}

	return chord, nil
}

// notify tells chord, potentialPred thinks it might be chord's predecessor.
func (chord *ChordServer) notify(potentialPred *chordrpc.Node) error {
	return nil
}

func (chord *ChordServer) join(joinNode *chordrpc.Node) error {
	if joinNode != nil {
		remote, err := chord.findSuccessorRPC(joinNode, chord.Id)
		if err != nil {
			return err
		}
		remotePred, err := chord.getPredecessorRPC(joinNode)
		if err != nil {
			return err
		}
		if idsEqual(remote.Id, chord.Id) {
			return errors.New("Duplicate IDs")
		}
		chord.fingerTable[0] = new(chordrpc.Node)
		chord.fingerTable[0].Id = remote.Id
		chord.fingerTable[0].Ip = remote.Ip
		chord.predecessor = new(chordrpc.Node)
		chord.predecessor.Id = remotePred.Id
		chord.predecessor.Ip = remotePred.Ip
	} else { // chord is the only one in the ring
		for i := 0; i < numBits; i++ {
			chord.fingerTable[i] = new(chordrpc.Node)
			chord.fingerTable[i].Id = chord.Id
			chord.fingerTable[i].Ip = chord.Ip
		}
		chord.predecessor = new(chordrpc.Node)
		chord.predecessor.Id = chord.Id
		chord.predecessor.Ip = chord.Ip
	}
	return nil
}

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

func (chord *ChordServer) findClosestPrecedingNode(id []byte) *chordrpc.Node {
	for i := numBits - 1; i >= 0; i-- {
		if chord.fingerTable[i] != nil && between(chord.fingerTable[i].Id, chord.Id, id) {
			return chord.fingerTable[i]
		}
	}
	return chord.Node
}

func hash(ipAddr string) []byte {
	h := sha1.New()
	h.Write([]byte(ipAddr))

	idInt := big.NewInt(0)
	idInt.SetBytes(h.Sum(nil)) // Sum() returns []byte, convert it into BigInt

	maxVal := big.NewInt(0)
	maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil) // calculate 2^m
	idInt.Mod(idInt, maxVal)                            // mod id to make it to be [0, 2^m - 1]
	if idInt.Cmp(big.NewInt(0)) == 0 {
		return []byte{0}
	}
	return idInt.Bytes()
}

func (chord *ChordServer) String() string {
	str := chord.Node.String() + "\n"

	str += "Finger table:\n"
	str += "ith | start | successor\n"
	for i := 0; i < numBits; i++ {
		currID := new(big.Int).SetBytes(chord.Id)
		offset := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(i)), nil)
		maxVal := big.NewInt(0)
		maxVal.Exp(big.NewInt(2), big.NewInt(numBits), nil)
		start := new(big.Int).Add(currID, offset)
		start.Mod(start, maxVal)
		successor := chord.fingerTable[i].Id
		str += fmt.Sprintf("%d   | %d     | %d\n", i, start, successor)
	}
	str += fmt.Sprintf("Predecessor: %v", chord.predecessor.Id)
	return str
}
