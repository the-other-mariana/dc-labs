package main

import (
	"flag"
	"fmt"
	"net/http"
	//"time"
	"io"
	"io/ioutil"
	"log"
	"strings"
	//"encoding/xml"
)

func responseHandler(res http.ResponseWriter, req *http.Request) {

	bucketName := req.FormValue("bucket")
	url := fmt.Sprintf("https://%v.s3.amazonaws.com", bucketName)
	resp, err := http.Get(url)
	if err != nil {
		panic("Error at S3 connection: " + err.Error())
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error at data reading: " + err.Error())
	}

	io.WriteString(res, "RESPONSE")
	io.WriteString(res, "\nURL: " + url)
    io.WriteString(res, "\nbucket: "+req.FormValue("bucket"))
    io.WriteString(res, "\ndir: "+req.FormValue("dir"))
	io.WriteString(res, "\n" + strings.Split(string(data), "\n")[0])
}
func main() {

	var port = flag.Int("port", 9000, "Port number.")
	flag.Parse()
	fmt.Printf("Service on Port: %v\n", *port) // port is a pointer

	socket := fmt.Sprintf("localhost:%v", *port)

	http.HandleFunc("/example", func(res http.ResponseWriter, req *http.Request) {
        responseHandler(res, req)
    })
	
	err := http.ListenAndServe(socket, nil)
	if err != nil {
		log.Fatal(err)
	}
}
