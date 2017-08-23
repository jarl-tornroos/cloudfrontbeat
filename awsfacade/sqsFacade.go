package awsfacade

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"fmt"
)

// SqsFacade handle simple SQS queue functions
type SqsFacade struct {
	client *sqs.SQS
	url    *string
}

// NewSqs is a construct function for creating the object
// with auth and name of the queue as arguments
func NewSqs(auth *Auth, name string) (*SqsFacade, error) {

	client := sqs.New(auth.Session, auth.Config)

	// Create queue if it doesn't exist
	queueInput := &sqs.CreateQueueInput{
		QueueName: aws.String(name),
		Attributes: map[string]*string{
			"VisibilityTimeout": aws.String("300"),
		},
	}
	response, err := client.CreateQueue(queueInput)
	if err != nil {
		return nil, err
	}

	sqsFacade := &SqsFacade{
		client: client,
		url:    response.QueueUrl,
	}

	return sqsFacade, nil
}

// SendMessage to queue
func (q *SqsFacade) SendMessage(message *string) error {
	messageInput := &sqs.SendMessageInput{
		MessageBody:  message,
		QueueUrl:     q.url,
		DelaySeconds: aws.Int64(1),
	}
	_, err := q.client.SendMessage(messageInput)

	return err
}

// ReceiveMessage from queue
func (q *SqsFacade) ReceiveMessage() (*sqs.Message, error) {
	messageInput := &sqs.ReceiveMessageInput{
		QueueUrl:            q.url,
		MaxNumberOfMessages: aws.Int64(1),
	}
	resp, err := q.client.ReceiveMessage(messageInput)
	if err != nil {
		return nil, err
	}
	// SQS messages should not be more than one
	if len(resp.Messages) > 1 {
		return nil, fmt.Errorf("Too many SQS messages")
	} else if len(resp.Messages) == 1 {
		return resp.Messages[0], nil
	} else {
		return &sqs.Message{}, nil
	}
}

// DeleteMessage from queue
func (q *SqsFacade) DeleteMessage(receiptHandle *string) error {
	var err error
	if receiptHandle != nil {
		messageInput := &sqs.DeleteMessageInput{
			QueueUrl:      q.url,
			ReceiptHandle: receiptHandle,
		}

		_, err = q.client.DeleteMessage(messageInput)
	}

	return err
}
