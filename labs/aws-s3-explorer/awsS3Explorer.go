package main

import (
	"flag"
	"fmt"
	"net/http"
	//"time"
	//"io"
	"io/ioutil"
	"log"
	"strings"
	"encoding/xml"
	"encoding/json"
)

type Content struct {
	Key string `xml:"Key"`
}

type XMLResult struct {
	XMLName  xml.Name  `xml:"ListBucketResult"`
	Name  string  `xml:"Name"`
	Contents []Content `xml:"Contents"`
}

type BucketDetails struct {
	BucketName string
	ObjectsCount int
	DirectoriesCount int
	Extensions map[string]int
}

func responseHandler(res http.ResponseWriter, req *http.Request) {

	extensions := make(map[string]int)
	directories := make(map[string]bool)
	objects := make(map[string]bool)

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

	for _,c := range xmlResult.Contents {

		objKey := c.Key

		if strings.HasSuffix(objKey, "/"){
			if !directories[objKey] {
				directories[objKey] = true
			}
		}

		if strings.Contains(objKey, ".") {

			_, exists := objects[objKey]

			if !exists {
				objects[objKey] = true
			}

			ext := strings.Split(objKey, ".")
			_, exists = extensions[ext[len(ext) - 1]];

			if !exists {
				extensions[ext[len(ext) - 1]] = 1
			}
			if exists {
				extensions[ext[len(ext) - 1]] += 1
			}
			
		}
	}

	bd := BucketDetails{
		BucketName: xmlResult.Name,
		ObjectsCount: len(objects),
		DirectoriesCount: len(directories),
		Extensions: extensions,
	}
	out, err := json.MarshalIndent(bd,"", "\t")
	if err != nil {
		panic(err)
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(out)
}
func main() {

	var port = flag.Int("port", 9000, "Port number.")
	flag.Parse()
	fmt.Printf("===== Service on port: %v =====\n", *port) // port is a pointer

	socket := fmt.Sprintf("localhost:%v", *port)

	http.HandleFunc("/example", func(res http.ResponseWriter, req *http.Request) {
        responseHandler(res, req)
    })
	
	err := http.ListenAndServe(socket, nil)
	if err != nil {
		log.Fatal(err)
	}
}
