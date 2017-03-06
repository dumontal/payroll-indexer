package main

import "fmt"

const (
	esURL   string = "http://127.0.0.1:9200"
	esIndex string = "payrolls"

	sqsURL    string = "https://sqs.eu-west-1.amazonaws.com/560569522348/payrolls"
	awsRegion string = "eu-west-1"
)

func main() {
	sqsClient := NewSQSClient(sqsURL, awsRegion)
	messages, err := sqsClient.Read()
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	if len(messages) == 0 {
		fmt.Println("Received no messages")
		return
	}

	//fmt.Println(messages)

	esClient := NewESClient(esURL, esIndex)
	esClient.EnsureIndexExists()
}
