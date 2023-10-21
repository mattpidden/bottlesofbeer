package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"time"
	"uk.ac.bris.cs/distributed2/bottles/stubs"

	//	"net/rpc"
	//	"fmt"
	//	"time"
	//	"net"
)

//RUNNING INSTRUCTIONS (for 3 instances)
// go run bottlesofbeer.go -this "8030" -next "localhost:8050"
// go run bottlesofbeer.go -this "8050" -next "localhost:8080"
// go run bottlesofbeer.go -this "8080" -next "localhost:8030" -n 99

var nextAddr string

func Singing(numberofbottles int) {
	stringBottles := strconv.Itoa(numberofbottles)
	response := ""
	if (numberofbottles == 1) {
		response = stringBottles + " bottles of beer on the wall, " + stringBottles + " bottles of beer. Take one down, and there's none left..."
	} else {
		response = stringBottles + " bottles of beer on the wall, " + stringBottles + " bottles of beer. Take one down, and pass it on..."
	}
	fmt.Println(response)
}

type BeerOperations struct {
	ResultChannel chan int
}

func (s *BeerOperations) SingLine(req stubs.Request, res *stubs.Response) (err error) {
	Singing(req.Bottles)
	res.Bottles = req.Bottles - 1
	s.ResultChannel <- res.Bottles
	return
}

func main(){
	//make channel to get bottles of beer from incoming request
	resultChannel := make(chan int)
	rpc.Register(&BeerOperations{ResultChannel: resultChannel})

	//Set up all flags for command line input
	thisPort := flag.String("this", "8030", "Port for this process to listen on")
	flag.StringVar(&nextAddr, "next", "localhost:8040", "IP:Port string for next member of the round.")
	bottles := flag.Int("n",0, "Bottles of Beer (launches song if not 0)")
	flag.Parse()

	//set up listening on this port
	listener, _ := net.Listen("tcp", ":"+*thisPort)
	defer listener.Close()

	//sleep to allow all instances to start before sending requests
	fmt.Println("Set up Listener - Sleeping 10 seconds")
	time.Sleep(10 * time.Second)

	//set up connection with next process
	client, _ := rpc.Dial("tcp", nextAddr)
	defer client.Close()

	//check if process given starting number of bottles
	if (*bottles != 0) {
		fmt.Println("This is the starting instance")

		//then tell next process to sing a line
		request := stubs.Request{Bottles: *bottles}
		response := new(stubs.Response)
		client.Call(stubs.SingLine, request, response)
	}

	//Create channel to let program finish execution
	doneChan := make(chan int)

	go func() {
		for {
			//Wait for instance to indicate that request has been heard and served
			leftoverBottles := <- resultChannel

			//Check if song is over
			if leftoverBottles == 0 {
				doneChan <- 1
			}

			// Give the song rhythm
			time.Sleep(100 * time.Millisecond)

			//Tell the next instance to sing a line
			request := stubs.Request{Bottles: leftoverBottles}
			response := new(stubs.Response)
			client.Call(stubs.SingLine, request, response)
		}
	}()

	//Concurrently listen for requests
	go func() {
		rpc.Accept(listener)
	}()

	//Blocks until song is done, then allow main to finish execution
	<- doneChan
}
