package chord

import (
	"crypto/sha1"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"time"
)

// import "fmt"
// import "os"

// Debug is enabled if set to 1
const Debug = 1

// DPrintf can be used for debug printing
func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}

// checkError checks err, prints err, aborts the program
func checkError(source string, err error) {
	if err != nil {
		log.Fatal(source+": ", err)
	}
}

// checkError checks err, prints err, this doesn't abort the program
func checkErrorGrace(source string, err error) {
	if err != nil {
		log.Println(source, ": ", err)
	}
}

func Hash(ipAddr string) []byte {
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

func ipGenerator(ip string, amount int) []string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	uniqueAddr := make(map[string]bool)

	addrs := make([]string, amount)

	for i := 0; i < amount; {
		ip := ip + ":" + strconv.Itoa(r1.Intn(10000)+10000)
		_, ok := uniqueAddr[ip]
		if !ok {
			uniqueAddr[ip] = true
			addrs[i] = ip
			i++
		}
	}

	return addrs
}
