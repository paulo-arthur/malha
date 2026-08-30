package main

import (
	"fmt"
	"io"
	"net/http"
)

func httpGETreq (url string, c chan int){
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
		var results_array []int 

		for i := 0; i < N; i++ {
			go httpGETreq(URL, c)
		}

		for i := 0; i < N; i++ {
			results_array = append(results_array, <-c)
		} 

		
	}
}