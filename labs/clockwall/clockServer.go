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

func getLocalTime(current time.Time, location string) time.Time {
	loc, error := time.LoadLocation(location)
	if error != nil { 
        panic(error) 
    }
	return current.In(loc)
}

func getResponse(sChan chan string, location string) {
	for {
		timeResponse := getLocalTime(time.Now(), location).Format("15:04:05\n")
		time.Sleep(1 * time.Second)
		sChan <- timeResponse
	}
}

func handleConn(c net.Conn, location string) {
	defer c.Close()

	sChan := make(chan string)
	go getResponse(sChan, location)

	for t := range sChan {
		_, err := io.WriteString(c, t)
		if err != nil {
			return 
		}
	}
}

func main() {
	timezone := os.Getenv("TZ")
	fmt.Printf("Timezone: %v\n", timezone)

	var port = flag.Int("port", 9000, "Port number.")
	flag.Parse()

	socket := fmt.Sprintf("localhost:%v", *port)
	listener, err := net.Listen("tcp", socket)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) 
			continue
		}
		handleConn(conn, timezone) 
	}
}
