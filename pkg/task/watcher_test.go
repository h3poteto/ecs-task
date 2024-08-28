package task

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	logstypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type mockedWatcher struct {
	LogsClient
	StreamsResp cloudwatchlogs.DescribeLogStreamsOutput
}

func (m mockedWatcher) DescribeLogStreams(ctx context.Context, params *cloudwatchlogs.DescribeLogStreamsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	return &m.StreamsResp, nil
}

func TestWaitStream(t *testing.T) {
	ctx := context.Background()
	oneOutput := cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []logstypes.LogStream{
			{
				Arn:           aws.String("Arn"),
				LogStreamName: aws.String("StreamName"),
			},
		},
	}
	one := &Watcher{
		awsLogs: mockedWatcher{StreamsResp: oneOutput},
		Group:   "Group",
		Stream:  "Stream",
	}
	stream, err := one.WaitStream(ctx)
	if err != nil {
		t.Error("Error waiting stream")
	}
	if *stream.LogStreamName != "StreamName" {
		t.Error("Stream name is invalid")
	}

	twoOutput := cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []logstypes.LogStream{
			{
				Arn:           aws.String("Arn1"),
				LogStreamName: aws.String("StreamName1"),
			},
			{
				Arn:           aws.String("Arn2"),
				LogStreamName: aws.String("StreamName2"),
			},
		},
	}
	two := &Watcher{
		awsLogs: mockedWatcher{StreamsResp: twoOutput},
		Group:   "Group",
		Stream:  "Stream",
	}
	_, err = two.WaitStream(ctx)
	if err == nil {
		t.Error("Does not error when multiple streams")
	}
}
