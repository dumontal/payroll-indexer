package main

import "log"

const (
	esURL   string = "http://127.0.0.1:9200"
	esIndex string = "payrolls"

	sqsURL    string = "https://sqs.eu-west-1.amazonaws.com/560569522348/payrolls"
	awsRegion string = "eu-west-1"

	nMessages int = 1000
)

func main() {
	sqsClient := NewSQSClient(sqsURL, awsRegion)
	esClient := NewESClient(esURL, esIndex)
	err := esClient.EnsureIndexExists()
	if err != nil {
		log.Println("Error", err)
		return
	}

	for i := 0; i < nMessages; i++ {
		var sqsMessages []SQSMessage
		sqsMessages, err = sqsClient.Read()
		if err != nil {
			log.Println("Error", err)
			return
		}

		for _, sqsMessage := range sqsMessages {
			err = esClient.Put(sqsMessage)
			if err != nil {
				log.Println("Error", err)
				return
			}
		}
	}

	err = esClient.Flush()
	if err != nil {
		log.Println("Error", err)
		return
	}

	log.Println("Done")
}
