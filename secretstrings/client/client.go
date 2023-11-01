package main

import (
	"bufio"
	"uk.ac.bris.cs/distributed2/secretstrings/stubs"

	//	"net/rpc"
	"flag"
	"log"
	"net/rpc"
	"os"
	//	"bufio"
	//	"os"
	//	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
	"fmt"
)

func main(){
	//Dealing with server and RPC
	server := flag.String("server","127.0.0.1:8030","IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	client, _ := rpc.Dial("tcp", *server)
	defer client.Close()

	//Dealing with wordlist file
	file, err := os.Open("wordlist")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		request := stubs.Request{scanner.Text()}
		response := new(stubs.Response)
		client.Call(stubs.PremiumReverseHandler, request, response)
		fmt.Println("Responded: " +response.Message)
	}
}
