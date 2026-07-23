package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Config struct {
	QueueURL string
}

type Message struct {
	Body          string
	ReceiptHandle string
}

type Consumer struct {
	client   *sqs.Client
	queueURL string
}

func NewConsumer(ctx context.Context, cfg Config) (*Consumer, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:   sqs.NewFromConfig(awsCfg),
		queueURL: cfg.QueueURL,
	}, nil
}

func (c *Consumer) Receive(ctx context.Context) ([]Message, error) {
	out, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]Message, 0, len(out.Messages))
	for _, msg := range out.Messages {
		messages = append(messages, Message{
			Body:          aws.ToString(msg.Body),
			ReceiptHandle: aws.ToString(msg.ReceiptHandle),
		})
	}

	return messages, nil
}

func (c *Consumer) Delete(ctx context.Context, receiptHandle string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})

	return err
}
