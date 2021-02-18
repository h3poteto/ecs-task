package task

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func TestBuildLogStream(t *testing.T) {
	tests := []struct {
		name     string
		arn      string
		expected string
	}{
		{
			name:     "TaskArnWithoutCluster",
			arn:      "arn:aws:ecs:ap-northeast-1:1234567890:task/c5cba4eb-5dad-405e-96db-71ef8eefe6a8",
			expected: "c5cba4eb-5dad-405e-96db-71ef8eefe6a8",
		},
		{
			name:     "TaskArnWithCluster",
			arn:      "arn:aws:ecs:ap-northeast-1:1234567890:task/my-cluster/c5cba4eb-5dad-405e-96db-71ef8eefe6a8",
			expected: "c5cba4eb-5dad-405e-96db-71ef8eefe6a8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{}
			ECSTask := &ecs.Task{
				TaskArn: aws.String(tt.arn),
			}
			taskID := task.buildLogStream(ECSTask)
			if taskID != tt.expected {
				t.Error("Task ID is invalid")
			}
		})
	}
}
