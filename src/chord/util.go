package chord

import "log"

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

func checkError(source string, err error) {
	if err != nil {
		log.Fatal(source, ": ", err)
	}
}
