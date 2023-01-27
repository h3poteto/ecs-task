package task

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Watcher has log group information and CloudWatchLogs Client.
type Watcher struct {
	awsLogs         cloudwatchlogsiface.CloudWatchLogsAPI
	Group           string
	Stream          string
	timestampFormat string
}

// NewWatcher returns a Watcher struct.
func NewWatcher(group, stream string, awsLogs cloudwatchlogsiface.CloudWatchLogsAPI, timestampFormat string) *Watcher {
	return &Watcher{
		Group:           group,
		Stream:          stream,
		awsLogs:         awsLogs,
		timestampFormat: timestampFormat,
	}
}

// GetStreams get cloudwatch logs streams according to log group name and stream prefix.
func (w *Watcher) GetStreams(ctx context.Context) ([]*cloudwatchlogs.LogStream, error) {
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(w.Group),
		LogStreamNamePrefix: aws.String(w.Stream),
		Descending:          aws.Bool(true),
	}
	output, err := w.awsLogs.DescribeLogStreamsWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return output.LogStreams, nil
}

// WaitStream waits until the log stream is generated.
func (w *Watcher) WaitStream(ctx context.Context) (*cloudwatchlogs.LogStream, error) {
	for {
		select {
		case <-time.After(2 * time.Second):
			streams, err := w.GetStreams(ctx)
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() == "Throttling" {
					log.Warn("Throttling")
					time.Sleep(5 * time.Second)
					continue
				}
			}
			if err != nil {
				return nil, err
			}
			if len(streams) == 1 {
				return streams[0], nil
			}
			if len(streams) > 1 {
				return nil, errors.New("There are multiple streams")
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Polling get log stream and print the logs with streaming.
func (w *Watcher) Polling(ctx context.Context) error {
	stream, err := w.WaitStream(ctx)
	if err != nil {
		return err
	}
	log.Infof("Log Stream: %+v", stream)
	fmt.Printf("Watching log stream: %s\n", *stream.Arn)
	var nextToken *string
	for {
		select {
		case <-time.After(2 * time.Second):
			input := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(w.Group),
				LogStreamName: stream.LogStreamName,
				StartFromHead: aws.Bool(true),
				NextToken:     nextToken,
			}
			output, err := w.awsLogs.GetLogEventsWithContext(ctx, input)
			if err != nil {
				return err
			}
			// Update next token
			nextToken = output.NextForwardToken
			w.printEvents(output.Events)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (w *Watcher) printEvents(events []*cloudwatchlogs.OutputLogEvent) {
	for _, event := range events {
		// AWS returns milliseconds of unix time.
		// So we have to transfer to second, nanoseconds.
		timestamp := time.Unix(*event.Timestamp / 1000, *event.Timestamp % 1000 * 1000000)
		message := *event.Message
		sTimestamp := timestamp.Format(w.timestampFormat)
		if sTimestamp != "" {
			sTimestamp += " "
		}
		fmt.Printf("%s%s\n", sTimestamp, message)
	}
}
