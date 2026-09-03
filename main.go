package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"net"
)

func readSinglePayloadFromFile(filePath string) string {
	fileContent, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Printf("Failed to read the file: %s", err)
	}

	return string(fileContent)
}

func readMultiplePayloadsFromFile(filePath string) ([]string, int) {
	rawFileContent := readSinglePayloadFromFile(filePath)
	var parsedObjectsList []map[string]interface{}

	err := json.Unmarshal([]byte(rawFileContent), &parsedObjectsList)
	if err != nil {
		fmt.Println("Error converting multiple payload into a slice.", err)
		return []string{}, 0
	}

	var formattedPayloadsList []string
	for _, singleObject := range parsedObjectsList {
		jsonBytes, _ := json.Marshal(singleObject)
		formattedPayloadsList = append(formattedPayloadsList, string(jsonBytes))
	}

	fmt.Println(formattedPayloadsList)
	return formattedPayloadsList, len(formattedPayloadsList)
}

func sendRawTCPRequest(url string, ip string, payload []string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("Something went wrong on connecting via TCP... > ", err)
	}

	defer conn.Close()

	reader := bufio.NewReader(conn)


	for i, p := range payload {
		_, err := conn.Write([]byte(p))
		if err != nil {
			log.Printf("Failed to send data: %v", err)
			return 
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read response from server: %v", err)
			return
		}

		fmt.Printf("Received from server: %s", response)

	}
}

func sendGetRequest(targetUrl string, statusCodeChannel chan int) {
	response, err := http.Get(targetUrl)
	if err != nil {
		fmt.Println("Request failed: ", err)
		statusCodeChannel <- 0
		return
	}

	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Failed to read body", err)
		statusCodeChannel <- 0
		return
	}

	fmt.Println("Status code: ", response.StatusCode)
	fmt.Println(string(responseBody))
	statusCodeChannel <- response.StatusCode
}

func sendPostRequest(targetUrl string, stringPayload string, statusCodeChannel chan int) {
	requestStructure, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer([]byte(stringPayload)))
	if err != nil {
		fmt.Print("Something went wrong...", err)
		return
	}

	//requestStructure.Header.Set("")

	httpClient := &http.Client{}
	response, err := httpClient.Do(requestStructure)
	if err != nil {
		fmt.Println("Something went wrong...", err)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)

	fmt.Println("Status Code: ", response.StatusCode)
	fmt.Println(string(responseBody))
}

func requestWorker(targetUrl string, jobsChannel <-chan string, resultsChannel chan<- int) {
	for currentJobPayload := range jobsChannel {
		sendPostRequest(targetUrl, currentJobPayload, resultsChannel)
	}
}

func main() {
	fmt.Println("-- MALIA - Go Fuzzler --")

	fmt.Printf("Enter the URL: ")
	var targetUrl string
	fmt.Scan(&targetUrl)

	fmt.Printf("Enter the method: ")
	var httpMethod string
	fmt.Scan(&httpMethod)

	fmt.Printf("Enter the number of simultaneous reqs: ")
	var workerCount int
	fmt.Scan(&workerCount)

	switch httpMethod {
	case "TCP":

	case "GET":
		statusCodeChannel := make(chan int)
		var statusCodesList []int

		for workerIndex := 0; workerIndex < workerCount; workerIndex++ {
			go sendGetRequest(targetUrl, statusCodeChannel)
		}

		for workerIndex := 0; workerIndex < workerCount; workerIndex++ {
			statusCodesList = append(statusCodesList, <-statusCodeChannel)
		}

	case "POST":
		fmt.Println("Choose the payload type: [L]ist | [U]nic")
		var attackType string
		fmt.Scan(&attackType)

		fmt.Printf("Enter the path to the payload file: ")
		var payloadFilePath string
		fmt.Scan(&payloadFilePath)

		switch attackType {
		case "L":
			jsonPayloadsList, totalPayloadsCount := readMultiplePayloadsFromFile(payloadFilePath)
			if totalPayloadsCount < workerCount {
				fmt.Println("More channels then requests detected. Reducing the number of channels to ", totalPayloadsCount)
				workerCount = totalPayloadsCount
			}

			jobsChannel := make(chan string, totalPayloadsCount)
			resultsChannel := make(chan int, totalPayloadsCount)

			for _, singlePayload := range jsonPayloadsList {
				jobsChannel <- singlePayload
			}
			close(jobsChannel)

			for workerIndex := 0; workerIndex < workerCount; workerIndex++ {
				go requestWorker(targetUrl, jobsChannel, resultsChannel)
			}

			for responseIndex := 0; responseIndex < len(payloads); responseIndex++ {
				receivedStatus := <-resultsChannel
				fmt.Println(receivedStatus)
			}

			fmt.Println(jsonPayload)
		case "U":
			statusCodeChannel := make(chan int)
			rawPayloadContent := readSinglePayloadFromFile(payloadFilePath)
			sendPostRequest(targetUrl, []byte(rawPayloadContent), statusCodeChannel)
		}
	}
}