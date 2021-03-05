// Clock Server is a concurrent TCP server that periodically writes the time.
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"flag"
	"os"

)

// function that transforms current time into another timezone time
func getLocalTime(current time.Time, location string) time.Time {
	loc, error := time.LoadLocation(location)
	if error != nil { 
        panic(error) 
    }
	return current.In(loc)
}

// function that waits 1 sec and sends a local time to the unbuffered channel
func getResponse(sChan chan string, location string) {
	for {
		timeResponse := getLocalTime(time.Now(), location).Format("15:04:05\n")
		time.Sleep(1 * time.Second)
		sChan <- timeResponse
	}
}

// function that calls a goroutine to fill the response channel
func handleConn(c net.Conn, location string) {
	defer c.Close()

	// UNBUFFERED CHANNEL FOR RESPONSE GATHERING
	sChan := make(chan string)
	go getResponse(sChan, location)

	// response sending
	for t := range sChan {
		response := fmt.Sprintf("%v" +" \t: " + "%v", location, t)
		_, err := io.WriteString(c, response)
		if err != nil {
			return 
		}
	}
}

// main goroutine that listens to an input port to serve it
func main() {
	timezone := os.Getenv("TZ")

	var port = flag.Int("port", 8030, "Port number.")
	flag.Parse()

	socket := fmt.Sprintf("localhost:%v", *port)
	listener, err := net.Listen("tcp", socket)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("=== Started clock server [%v] at port [%v] ====\n", timezone, *port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) 
			continue
		}
		handleConn(conn, timezone) 
	}
}
