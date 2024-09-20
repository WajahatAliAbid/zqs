package helper

import (
	"context"
	"encoding/hex"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
)

type SqsClient struct {
	client *sqs.Client
}

func New(
	cfg *aws.Config,
) *SqsClient {
	_client := sqs.NewFromConfig(*cfg)
	return &SqsClient{
		client: _client,
	}
}

type QueueInfo struct {
	Name       string             `json:"name"`
	Url        string             `json:"url"`
	Tags       *map[string]string `json:"tags"`
	Attributes *map[string]string `json:"attributes"`
}

func (s *SqsClient) GetQueueUrl(
	ctx context.Context,
	queueName *string,
) (*string, error) {
	output, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: queueName,
	})
	if err != nil {
		return nil, err
	}
	return output.QueueUrl, nil
}

func uuidv5(data string) string {
	bytes := []byte(data)
	id := uuid.NewSHA1(uuid.NameSpaceURL, bytes)
	return hex.EncodeToString(id[:])
}

func (s *SqsClient) SendMessages(
	ctx context.Context,
	queueUrl *string,
	messages *[]map[string]interface{},
) error {
	if len(*messages) == 0 {
		return nil
	}
	entries := []types.SendMessageBatchRequestEntry{}
	for _, message := range *messages {
		obj, _ := json.Marshal(message)
		str_obj := string(obj)
		id := uuidv5(str_obj)
		entries = append(entries, types.SendMessageBatchRequestEntry{
			Id:          &id,
			MessageBody: &str_obj,
		})
	}
	_, err := s.client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: queueUrl,
	})
	if err != nil {
		return err
	}
	return nil
}
