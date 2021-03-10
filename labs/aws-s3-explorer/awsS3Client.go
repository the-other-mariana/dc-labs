// client that requests to server (via proxy param) a bucket and directory (optional)
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"flag"
)

func main() {

	var proxy = flag.String("proxy", "localhost:8000", "Proxy url.")
	var bucketName = flag.String("bucket", "", "S3 bucket name.")
	var directory = flag.String("directory", "", "Directory name.")
	flag.Parse()
	
	if *bucketName == "" {
		fmt.Println("ERROR - Missing parameters.")
		return
	}

	request := fmt.Sprintf("http://%v/example?bucket=%v&dir=%v", *proxy, *bucketName, *directory)
	resp, err := http.Get(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
