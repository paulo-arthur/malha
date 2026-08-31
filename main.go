package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"bytes"
)

func readFile(path string) string {
	content, err := os.ReadFile(path)
	
	if err != nil {
		fmt.Printf("Failed to read the file: %s", err)
	}

	return string(content)
}

func httpGETreq(url string, c chan int){
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request failed: ", err)
		c <- 0
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read body", err)
		c <- 0
		return
	}

	fmt.Println("Status code: ", resp.StatusCode)
	fmt.Println(string(body))
	c <- resp.StatusCode
}

func httpPOSTreq(url string, filepath string, c chan int) {
	binary_payload := []byte(readFile(filepath))

	req_structure, err := http.NewRequest("POST", url, bytes.NewBuffer(binary_payload))
	if err != nil {
		fmt.Print("Something went wrong...", err)
		return 
	}

	//req_structure.Header.Set("")

	client := &http.Client{}
	resp, err := client.Do(req_structure)
	if err != nil {
		fmt.Println("Something went wrong...", err)
		return 
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	fmt.Println("Status Code: ", resp.StatusCode)
	fmt.Println(string(bodyBytes))
}

func main() {
	fmt.Println("-- MALIA - Go Fuzzler")

	fmt.Printf("Enter the URL: ")
	var URL string
	fmt.Scan(&URL)

	fmt.Printf("Enter the method: ")
	var method string
	fmt.Scan(&method)

	fmt.Printf("Enter the number of simultaneous reqs: ")
	var N int
	fmt.Scan(&N)

	switch method {
	case "GET":
		c := make(chan int)
		var results_array []int 

		for i := 0; i < N; i++ {
			go httpGETreq(URL, c)
		}

		for i := 0; i < N; i++ {
			results_array = append(results_array, <-c)
		} 

	case "POST":
		fmt.Printf("Enter the path to the payload file: ")
		var path string
		fmt.Scan(&path)
		
		c := make(chan int)
		httpPOSTreq(URL, path, c)

		
	}
}