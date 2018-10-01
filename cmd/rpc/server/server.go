package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	server "github.com/philipid3s/go-rpc/api/rpc/v1"
)

func main() {
	arith := new(server.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	fmt.Printf("Server started: %s", l.Addr().String())
	http.Serve(l, nil)
}
