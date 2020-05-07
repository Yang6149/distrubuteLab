package shardmaster

import "log"

// Debugging
const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}
func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)

}
