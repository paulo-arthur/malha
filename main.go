package main

import (
	"fmt"
	"io"
	"net/http"
)

func httpReq (url string, c chan int){
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
	var URL string = "https://webhook.site/82fb4d9c-7045-4de7-b019-444bcfecc16f"

	c := make(chan int)

	for i := 0; i < 50; i++ {
		go httpReq(URL, c)
	}

	for i := 0; i < 50; i++ {
		fmt.Println(<-c)
	}
}