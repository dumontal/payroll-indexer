package main

/*
  Amazon SQS package
*/

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSClient is the Amazon client we use in this POC.
type SQSClient struct {
	svc  *sqs.SQS
	qURL string
}

// SQSMessage is the structure of the messages we want to read.
type SQSMessage struct {
	CompanyID  string `json:"companyId"`
	EmployeeID string `json:"employeeId"`
	AbsMonth   int    `json:"absMonth"`
	Content    struct {
		S3Link string `json:"s3Link"`
		Text   string `json:"text"`
	} `json:"content"`
}

// NewSQSClient creates a new instance of the Amazon SQS client.
func NewSQSClient(qURL, region string) *SQSClient {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return &SQSClient{
		qURL: qURL,
		svc:  sqs.New(sess),
	}
}

// Read reads a new message from the SQS queue.
// Returns an array of read messages.
func (c *SQSClient) Read() ([]SQSMessage, error) {
	result, err := c.svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &c.qURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(0),
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		return nil, err
	}

	var messages []SQSMessage
	for _, sqsMessage := range result.Messages {
		var message SQSMessage
		err = json.Unmarshal([]byte(*(sqsMessage.Body)), &message)
		if err != nil {
			return messages, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}
