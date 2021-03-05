// code that acts as a read-only TCP client for several clock servers at once.
package main

import (
	"io"
	"log"
	"net"
	"os"
	//"fmt"
	"strings"
	//"time"
)

func dialServer(socket string, c chan int) {
	conn, err := net.Dial("tcp", socket)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	out := os.Stdout
	io.Copy(out, conn)
	c <- 1 // sends 1 to buffered channel that can receive up to 3 values before blocking
}

func main() {
	//printerChan := make(chan LocalTime)
	c := make(chan int, 3)
	args := os.Args[1:]
	var sockets []string

	for _,arg := range args {
		socket := (strings.Split(arg, "="))[1]
		sockets = append(sockets, socket)
	}

	for _,socket := range sockets {
		go dialServer(socket, c)
	}
	<- c //receiver, this line waits until c is full to proceed
	close(c)
	
}