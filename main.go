package main

import (
	"log"
	"server/server"
)

func main() {
	log.Println("Hello")
	//srv := server.NewTcpServer("127.0.0.1", "5555")
	srv := server.NewTcpServer("127.0.0.1", "5555")
	srv.StartServer()
}
