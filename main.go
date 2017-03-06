package main

import "fmt"

const (
	sqsURL    string = "https://sqs.eu-west-1.amazonaws.com/560569522348/payrolls"
	awsRegion string = "eu-west-1"
)

func main() {
	sqsClient := NewSQSClient(sqsURL, awsRegion)
	message, err := sqsClient.Read()
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	if message == nil {
		fmt.Println("Received no messages")
		return
	}

	fmt.Println(message)
}

/*
func pingES() {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Obtain a client and connect to the default Elasticsearch installation
	// on 127.0.0.1:9200. Of course you can configure your client to connect
	// to other hosts and configure it in various other ways.
	client, err := elastic.NewClient()
	if err != nil {
		// Handle error
		panic(err)
	}

	// Ping the Elasticsearch server to get e.g. the version number
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
}
*/
