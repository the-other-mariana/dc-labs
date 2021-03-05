// code that acts as a read-only TCP client for several clock servers at once.
package main

import (
	"io"
	"log"
	"net"
	"os"
	"strings"
	"strconv"
)

// function that connects to a clock server and gets the time
func dialServer(socket string, c chan int) {
	conn, err := net.Dial("tcp", socket)
	if err != nil {
		log.Fatal(err)
	}

	port, err := strconv.Atoi((strings.Split(socket, ":"))[1])
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()
	_, err = io.Copy(os.Stdout, conn)
	if err != nil {
		log.Fatal(err)
	}

	c <- port
	// sends 1 to buffered channel that can receive up to 3 values before blocking
}

// main goroutine
func main() {
	
	args := os.Args[1:]
	var sockets []string

	for _,arg := range args {
		socket := (strings.Split(arg, "="))[1]
		sockets = append(sockets, socket)
	}

	// BUFFERED CHANNEL WITH CAPACITY AS LARGE AS NUMBER OF PORTS
	c := make(chan int, len(sockets)) 

	for _,socket := range sockets {
		// calling a separate goroutine for connecting to each of the sockets
		go dialServer(socket, c)
	}

	//receiver, this line makes the program wait until c has a value to proceed
	<- c 
	close(c)
	
}