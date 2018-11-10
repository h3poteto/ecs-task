package task

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func TestBuildLogStream(t *testing.T) {
	task := &Task{}
	ECSTask := &ecs.Task{
		TaskArn: aws.String("arn:aws:ecs:ap-northeast-1:1234567890:task/c5cba4eb-5dad-405e-96db-71ef8eefe6a8"),
	}
	taskID := task.buildLogStream(ECSTask)
	if taskID != "c5cba4eb-5dad-405e-96db-71ef8eefe6a8" {
		t.Error("Task ID is invalid")
	}
}
