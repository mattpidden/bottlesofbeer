
package main
import (
//	"errors"
//	"flag"
//	"fmt"
//	"net"
	"time"
	"math/rand"
//	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
//	"net/rpc"
)

/** Super-Secret `reversing a string' method we can't allow clients to see. **/
func ReverseString(s string, i int) string {
    time.Sleep(time.Duration(rand.Intn(i))* time.Second)
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

