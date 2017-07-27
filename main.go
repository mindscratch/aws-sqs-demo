package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	foo := ec2metadata.New(sess)
	region, err := foo.Region()
	if err != nil {
		fmt.Println("error getting region", err)
		return
	}
	fmt.Println("REGION:", region)
	sess.Config.Region = &region

	// Create a SQS service client.
	svc := sqs.New(sess)

	result, err := svc.ListQueues(nil)
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("Success")
	// As these are pointers, printing them out directly would not be useful.
	for i, urls := range result.QueueUrls {
		// Avoid dereferencing a nil pointer.
		if urls == nil {
			continue
		}
		fmt.Printf("%d: %s\n", i, *urls)
	}

	// get the queue url
	result2, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String("test"),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	fmt.Println("QUEUE URL IS", *result2.QueueUrl)

	// SEND A MESSAGE
	qURL := result2.QueueUrl

	if 1 == 2 {
		sendResult, err := svc.SendMessage(&sqs.SendMessageInput{
			DelaySeconds: aws.Int64(0),
			MessageAttributes: map[string]*sqs.MessageAttributeValue{
				"Title": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String("The Whistler"),
				},
				"Author": &sqs.MessageAttributeValue{
					DataType:    aws.String("String"),
					StringValue: aws.String("John Grisham"),
				},
				"WeeksOn": &sqs.MessageAttributeValue{
					DataType:    aws.String("Number"),
					StringValue: aws.String("6"),
				},
			},
			MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
			QueueUrl:    qURL,
		})

		if err != nil {
			fmt.Println("Error", err)
			return
		}

		fmt.Println("Success", *sendResult.MessageId)
	}
	// READ THE MESSAGE
	readResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            qURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(0),
		WaitTimeSeconds:     aws.Int64(0),
	})

	if err != nil {
		fmt.Println("Error", err)
		return
	}

	if len(readResult.Messages) == 0 {
		fmt.Println("Received no messages")
		return
	}

	for _, msg := range readResult.Messages {
		fmt.Printf("READ [%s] %s\n", *msg.MessageId, *msg.Body)

		// DELETE THE MESSAGE
		resultDelete, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      qURL,
			ReceiptHandle: msg.ReceiptHandle,
		})

		if err != nil {
			fmt.Println("Delete Error", err)
			return
		}

		fmt.Println("Message Deleted", resultDelete)
	}
}
