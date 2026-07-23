package queue

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Producer struct {
	client *sqs.Client
	queueURL string
}

func NewProducer(ctx context.Context, cfg Config) (*Producer, error) {
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &Producer{
		client: sqs.NewFromConfig(awsCfg),
		queueURL: cfg.QueueURL,
	}, nil
}

func (p *Producer) Send(ctx context.Context, body string) error {
	_, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl: aws.String(p.queueURL),
		MessageBody: aws.String(body),
	})
	return  err
}

func (p *Producer) SendChunk(ctx context.Context, bodies []string) error {
	entries := make([]types.SendMessageBatchRequestEntry, len(bodies))
	for i, body := range bodies {
		entries[i] = types.SendMessageBatchRequestEntry{
			Id: aws.String(strconv.Itoa(i)),
			MessageBody: aws.String(body),
		}
	}

	out, err := p.client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(p.queueURL),
		Entries: entries,
	})
	if err != nil {
		return err
	}
	if len(out.Failed) > 0 {
		return fmt.Errorf("%d of %d messages failed to send", len(out.Failed), len(bodies))
	}

	return  nil
}

func (p *Producer) SendBatch(ctx context.Context, bodies []string) error {
	for start := 0; start < len(bodies); start += 10 {
		end := min(start+10, len(bodies))
		if err := p.SendChunk(ctx, bodies[start:end]); err != nil {
			return err
		}
	}

	return  nil
}