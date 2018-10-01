package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"

	server "github.com/philipid3s/go-rpc/api/rpc/v1"
)

func main() {

	if len(os.Args) < 3 {
		log.Fatal("client call error: 2 arguments required")
	}

	var err error

	nums := make([]int, len(os.Args))

	for i := 1; i < len(os.Args); i++ {
		if nums[i], err = strconv.Atoi(os.Args[i]); err != nil {
			log.Fatal("converting args:", err)
		}
	}

	// RPC server address
	address := ""

	client, err := rpc.DialHTTP("tcp", address+":1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	args := &server.Args{nums[1], nums[2]}

	// Asynchronous call
	quotient := new(server.Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	if replyCall.Error != nil {
		log.Fatal("Arith error:", err)
	}
	fmt.Printf("Divide: %d/%d=%d\n", args.A, args.B, quotient.Quo)

	// Synchronous call
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("Arith error:", err)
	}
	fmt.Printf("Multiply: %d*%d=%d\n", args.A, args.B, reply)
}
