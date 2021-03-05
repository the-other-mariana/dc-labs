// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"io"
	"log"
	"net"
	"time"
)

func getResponse(sChan chan string) {
	for {
		timeResponse := time.Now().Format("15:04:05\n")
		sChan <- timeResponse
	}
}

func handleConn(c net.Conn, sChan chan string) {
	defer c.Close()
	go getResponse(sChan)
	for t := range sChan {
		_, err := io.WriteString(c, t)
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	sChan := make(chan string)
	listener, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) 
			continue
		}
		go handleConn(conn, sChan) 
	}
}
