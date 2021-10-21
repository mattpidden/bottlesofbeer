package main
import (
//	"net/rpc"
	"flag"
//	"bufio"
//	"os"
//	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
	"fmt"
)

func main(){
	server := flag.String("server","127.0.0.1:8030","IP:port string to connect to as server")
	flag.Parse()
	fmt.Println("Server: ", *server)
	//TODO: connect to the RPC server and send the request(s)
}
