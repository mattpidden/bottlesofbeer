
package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"

	//	"errors"
	//	"flag"
	//	"fmt"
	//	"net"
	"time"
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

type SecretStringOperations struct {
	ResultChannel chan string
}

func (s *SecretStringOperations) Reverse(req stubs.Request, res *stubs.Response) (err error) {
	res.Message = ReverseString(req.Message, 10)
	return
}

func (s *SecretStringOperations) FastReverse(req stubs.Request, res *stubs.Response) (err error) {
	res.Message = ReverseString(req.Message, 2)
	s.ResultChannel <- req.Message
	return
}

func main() {
	pAddr := flag.String("port", "8030", "Port to listen on")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	resultChannel := make(chan string)
	rpc.Register(&SecretStringOperations{ResultChannel: resultChannel})

	go func() {
		for {
			result := <-resultChannel
			fmt.Println(result)
			// Do something with the result here
		}
	}()

	listener, _ := net.Listen("tcp", ":"+*pAddr)
	defer listener.Close()
	rpc.Accept(listener)
}