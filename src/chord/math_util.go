package chord

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
)

func intEqualByte(x int, y []byte) bool {
	xBytes := big.NewInt(int64(x)).Bytes()
	return idsEqual(xBytes, y)
}

// add takes one number in bytes and second number in int64, return the result in bytes
// does not work with negative numbers
func addBytesInt64(numberInBytes []byte, addend int64) []byte {
	addend1 := big.NewInt(0).SetBytes(numberInBytes)
	addend2 := big.NewInt(addend)
	return addend1.Add(addend1, addend2).Bytes()
}

func addIDs(x, y []byte) []byte {
	if len(x) == 0 || len(y) == 0 {
		panic("empty byte array")
	}
	xInt := big.NewInt(0).SetBytes(x)
	yInt := big.NewInt(0).SetBytes(y)

	xInt.Add(xInt, yInt)
	return xInt.Bytes()
}

func addBytesBigint(numberInBytes []byte, addend *big.Int) []byte {
	if len(numberInBytes) == 0 {
		panic("empty byte array")
	}
	addend1 := big.NewInt(0).SetBytes(numberInBytes)
	return addend.Add(addend, addend1).Bytes()
}

func idsEqual(x, y []byte) bool {
	if len(x) == 0 || len(y) == 0 {
		panic("idsEqual: empty byte array")
	}
	return bytes.Equal(x, y)
}

// Returns true if x is between lo and hi
func between(id []byte, lo []byte, hi []byte) bool {
	if len(id) == 0 || len(lo) == 0 || len(hi) == 0 {
		panic("between: empty byte array")
	}
	idInt := big.NewInt(0).SetBytes(id)
	loInt := big.NewInt(0).SetBytes(lo)
	hiInt := big.NewInt(0).SetBytes(hi)

	switch loInt.Cmp(hiInt) {
	case -1: // lo < hi
		return (idInt.Cmp(loInt) == 1) && (idInt.Cmp(hiInt) == -1)
	case 1: // lo > hi
		return (idInt.Cmp(loInt) > 0) || (idInt.Cmp(hiInt) < 0)
	case 0: // lo == hi
		return idInt.Cmp(loInt) != 0
	}

	// (2, 3)
	return false
}

// Returns true if begin <= target < end, in the ring
func betweenLeftInclusive(id []byte, lo []byte, hi []byte) bool {
	if len(id) == 0 || len(lo) == 0 || len(hi) == 0 {
		panic("empty byte array")
	}
	return between(id, lo, hi) || bytes.Equal(id, lo)
}

func betweenRightInclusive(id []byte, lo []byte, hi []byte) bool {
	if len(id) == 0 || len(lo) == 0 || len(hi) == 0 {
		panic("empty byte array")
	}
	return between(id, lo, hi) || bytes.Equal(id, hi)
}

func intToByteArray(i int) []byte {
	if i == 0 {
		return []byte{0}
	}
	return big.NewInt(int64(i)).Bytes()
}

func reverseHash(exp int, ipAddr string, startPortNum int) []string {
	numNodes := int(math.Pow(2, float64(exp)))
	ips := make([]string, numNodes)
	for i := 0; i < numNodes; startPortNum++ {
		testAddr := fmt.Sprintf("%s:%d", ipAddr, startPortNum)
		hashResult := Hash(testAddr)
		if idsEqual(intToByteArray(i), hashResult) {
			ips[i] = testAddr
			i++
		}
	}
	return ips
}
