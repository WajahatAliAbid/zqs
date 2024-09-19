package helper

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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

func stringInSlice(sub string, list []string, partialMatch bool) bool {
	for _, b := range list {
		if b == sub {
			return true
		}
		if partialMatch && strings.Contains(strings.ToLower(b), strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

func (s *SqsClient) GetQueueAttributes(
	ctx context.Context,
	queueUrl *string,
) (*map[string]string, error) {
	output, err := s.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl: queueUrl,
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
	})
	if err != nil {
		return nil, err
	}
	return &output.Attributes, nil
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

func (s *SqsClient) ListQueueTags(
	ctx context.Context,
	queueUrl *string,
) (*map[string]string, error) {
	output, err := s.client.ListQueueTags(ctx, &sqs.ListQueueTagsInput{
		QueueUrl: queueUrl,
	})
	if err != nil {
		return nil, err
	}
	return &output.Tags, nil
}

func (s *SqsClient) ListQueues(
	ctx context.Context,
	names *[]string,
	withAttributes bool,
	partialMatch bool,
) (*[]QueueInfo, error) {
	output, err := s.client.ListQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		return nil, err
	}

	queues := []QueueInfo{}

	for queue := range output.QueueUrls {
		// get name of queue by splitting url
		queueUrl := output.QueueUrls[queue]
		splits := strings.Split(queueUrl, "/")
		name := splits[len(splits)-1]

		if names != nil && !stringInSlice(name, *names, partialMatch) {
			continue
		}
		info := QueueInfo{
			Name: name,
			Url:  queueUrl,
		}

		if withAttributes {
			attributes, err := s.GetQueueAttributes(ctx, &queueUrl)
			if err != nil {
				return nil, err
			}
			info.Attributes = attributes

			tags, err := s.ListQueueTags(ctx, &queueUrl)
			if err != nil {
				return nil, err
			}

			info.Tags = tags
		}

		queues = append(queues, info)
	}

	return &queues, nil
}

func (s *SqsClient) ListQueuesWithAttributes(
	ctx context.Context,
	names *[]string,
) (*[]QueueInfo, error) {
	return s.ListQueues(ctx, names, true, true)
}
