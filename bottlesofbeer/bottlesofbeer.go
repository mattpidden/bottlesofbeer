package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"time"
	"uk.ac.bris.cs/distributed2/bottles/stubs"

	//	"net/rpc"
	//	"fmt"
	//	"time"
	//	"net"
)

var nextAddr string

func Singing(numberofbottles int) {
	if (numberofbottles == 1) {
		fmt.Println("%v bottles of beer on the wall, %v bottles of beer. Take one down, and there's none left...", numberofbottles)
	} else {
		fmt.Println("%v bottles of beer on the wall, %v bottles of beer. Take one down, pass it around...", numberofbottles)
	}
}

type BeerOperations struct {}

func (s *BeerOperations) SingLine(req stubs.Request) (err error) {
	Singing(req.Bottles)
	return
}

func main(){
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	rpc.Register(&BeerOperations{})

	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")

	bottles := flag.Int("n",0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()

	//set up listening on this port
	listener, _ := net.Listen("tcp", ":"+*thisPort)
	defer listener.Close()

	//set up connection with next process
	client, _ := rpc.Dial("tcp", nextAddr)
	defer client.Close()

	//sleep 20 seconds to allow all the process to begin their listening function
	time.Sleep(20 * time.Second)


	//chcek if process given starting number of bottles
	if (*bottles != 0) {
		//start the song here
		Singing(*bottles)

		//then tell next process to sing a line
		request := stubs.Request{*bottles}
		response := new(stubs.Response)
		client.Call(stubs.SingLine, request, response)
	}
	for {
		//listen for when it is your turn
		rpc.Accept(listener)

		//tell next process to sing a line
		request := stubs.Request{*bottles}
		response := new(stubs.Response)
		client.Call(stubs.SingLine, request, response)
	}
}
