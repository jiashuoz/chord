package kvchord

import (
	"log"
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
		log.Fatal(source, ": ", err)
	}
}

// checkError checks err, prints err, this doesn't abort the program
func checkErrorGrace(source string, err error) {
	if err != nil {
		log.Println(source, ": ", err)
	}
}
