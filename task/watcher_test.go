package task

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
)

type mockedWatcher struct {
	cloudwatchlogsiface.CloudWatchLogsAPI
	StreamsResp cloudwatchlogs.DescribeLogStreamsOutput
}

func (m mockedWatcher) DescribeLogStreams(in *cloudwatchlogs.DescribeLogStreamsInput) (*cloudwatchlogs.DescribeLogStreamsOutput, error) {
	return &m.StreamsResp, nil
}

func TestWaitStream(t *testing.T) {
	ctx := context.Background()
	oneOutput := cloudwatchlogs.DescribeLogStreamsOutput{
		LogStreams: []*cloudwatchlogs.LogStream{
			&cloudwatchlogs.LogStream{
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
		LogStreams: []*cloudwatchlogs.LogStream{
			&cloudwatchlogs.LogStream{
				Arn:           aws.String("Arn1"),
				LogStreamName: aws.String("StreamName1"),
			},
			&cloudwatchlogs.LogStream{
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
