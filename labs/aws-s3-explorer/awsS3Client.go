package main

import (
	"bufio"
	"fmt"
	"net/http"
	"flag"
)

func main() {

	var proxy = flag.String("proxy", "localhost:8080", "Proxy url.")
	var bucketName = flag.String("bucket", "none", "S3 bucket name.")
	var directory = flag.String("directory", "", "Directory name.")
	flag.Parse()
	fmt.Printf("Proxy: %v\n", *proxy)
	fmt.Printf("Bucket Name: %v\n", *bucketName)
	fmt.Printf("Directory: %v\n", *directory)

	request := fmt.Sprintf("http://%v/example?bucket=%v&dir=%v", *proxy, *bucketName, *directory)
	resp, err := http.Get(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Print the first 5 lines of the response body.
	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
