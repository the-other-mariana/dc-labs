// server that serves the bucket and directory details in JSON format
package main

import (
	"flag"
	"fmt"
	"net/http"
	"io/ioutil"
	"io"
	"log"
	"strings"
	"encoding/xml"
	"encoding/json"
)

// key objects from xml
type Content struct {
	Key string `xml:"Key"`
}

// xml data struct receiver
type XMLResult struct {
	XMLName  xml.Name  `xml:"ListBucketResult"`
	Name  string  `xml:"Name"`
	Contents []Content `xml:"Contents"`
}

// result struct 1
type BucketDetails struct {
	BucketName string
	ObjectsCount int
	DirectoriesCount int
	Extensions map[string]int
}

// result struct 2
type DirDetails struct {
	BucketName string
	DirectoryName string
	ObjectsCount int
	DirectoriesCount int
	Extensions map[string]int
}

func responseHandler(res http.ResponseWriter, req *http.Request) {

	extensions := make(map[string]int)
	directories := make(map[string]bool)
	objects := make(map[string]bool)

	bucketName := req.FormValue("bucket")
	reqDir := req.FormValue("dir")
	url := fmt.Sprintf("https://%v.s3.amazonaws.com", bucketName)

	// get the url response from S3
	resp, gerr := http.Get(url)
	if gerr != nil {
		println("ERROR - Get S3 connection error.")
		io.WriteString(res, "ERROR - Get S3 connection error.")
		return
	}
	defer resp.Body.Close()

	data, rerr := ioutil.ReadAll(resp.Body)
	if rerr != nil {
		println("ERROR - Http response reading error.")
		io.WriteString(res, "ERROR - Http response reading error.")
		return
	}

	// decode xml into struct
	var xmlResult XMLResult
	xerr := xml.Unmarshal(data, &xmlResult)
	if xerr != nil {
		println("ERROR - Permission denied bucket.")
		io.WriteString(res, "ERROR - Permission denied bucket.\n")
		return
	}

	// read and process struct content
	for _,c := range xmlResult.Contents {
		objKey := c.Key
		dir := fmt.Sprintf("%v/", reqDir)

		// conditions to filter directory cases
		if reqDir != "" && !strings.HasPrefix(objKey, dir){
			continue
		}

		if reqDir != "" && strings.HasPrefix(objKey, dir){
			objKey = strings.Replace(objKey, dir, "", 1)
			if objKey == ""{
				// remaining directory is the request directory itself
				continue
			}
		}

		// process the key string
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

	// create result struct with final dictionaries and transform it to JSON
	if reqDir == "" {
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
	if reqDir != "" {
		dd := DirDetails{
			BucketName: xmlResult.Name,
			DirectoryName: reqDir,
			ObjectsCount: len(objects),
			DirectoriesCount: len(directories),
			Extensions: extensions,
		}
	
		out, err := json.MarshalIndent(dd, "", "\t")
		if err != nil {
			panic(err)
		}
		res.Header().Set("Content-Type", "application/json")
		res.Write(out)
	}
}
func main() {

	var port = flag.Int("port", 9000, "Port number.")
	flag.Parse()
	fmt.Printf("===== Service on port: %v =====\n", *port) // port is a pointer

	socket := fmt.Sprintf("localhost:%v", *port)

	// get the url from cliente request
	http.HandleFunc("/example", func(res http.ResponseWriter, req *http.Request) {
        responseHandler(res, req)
    })
	
	err := http.ListenAndServe(socket, nil)
	if err != nil {
		log.Fatal(err)
	}
}
