package main

import (
	"fmt"
	"io"
	"net/http"
)

func httpGETreq (url string, c chan int){
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request failed.", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read body", err)
		return
	}

	fmt.Println("Status code: ", resp.StatusCode)
	fmt.Println(string(body))
	c <- resp.StatusCode
}

func main() {
	fmt.Println("Malia")

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
		var results_array []string 

		for _ := 0; _ < N; _++ {
			go httpGETreq(URL, c)
		}

		for _ := 0; _ < N; _++ {
			results_array = append(results_array, <-c)
		} 


	}
}