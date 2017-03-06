package main

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
// Returns nil if there is no message to read (empty SQS for instance).
func (c *SQSClient) Read() (*SQSMessage, error) {
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

	if len(result.Messages) == 0 {
		return nil, nil
	}

	var message SQSMessage
	toDecode := *(result.Messages[0].Body)
	err = json.Unmarshal([]byte(toDecode), &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
