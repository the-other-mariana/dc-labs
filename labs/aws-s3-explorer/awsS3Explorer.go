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
	"encoding/xml"
)

type Content struct {
	Key string `xml:"Key"`
}

type XMLResult struct {
	XMLName  xml.Name  `xml:"ListBucketResult"`
	Name  string  `xml:"Name"`
	Contents []Content `xml:"Contents"`
}

func responseHandler(res http.ResponseWriter, req *http.Request) {

	extensions := make(map[string]int)
	directories := make(map[string]bool)

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
	defer resp.Body.Close()

	var xmlResult XMLResult
	xml.Unmarshal(data, &xmlResult)

	strRes := ""
	strRes += fmt.Sprintf("Bucket name: %v\n", xmlResult.Name)
	
	for _,c := range xmlResult.Contents {
		objKey := c.Key
		if strings.HasSuffix(objKey, "/"){
			if !directories[objKey] {
				directories[objKey] = true
			}
		}
		if strings.Contains(objKey, ".") {
			ext := strings.Split(objKey, ".")
			extensions[ext[len(ext) - 1]] += 1
		}
		//strRes += fmt.Sprintf("Content %v: Key: %v\n", index, c.Key)
	}

	strRes += fmt.Sprintf("Directories: %v\n", len(directories))
	strRes += fmt.Sprintf("Extensions: %v\n", len(extensions))
	for key, value := range extensions {
		strRes += fmt.Sprintf("%v: %v\n", key, value)
	}

	io.WriteString(res, "RESPONSE")
	io.WriteString(res, "\nURL: " + url)
	io.WriteString(res, "\n" + strRes)
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
