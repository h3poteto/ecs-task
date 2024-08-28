package task

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockedDescribeTaskDefinition struct {
	TaskDefinitionClient
	Resp ecs.DescribeTaskDefinitionOutput
}

func (m mockedDescribeTaskDefinition) DescribeTaskDefinition(ctx context.Context, params *ecs.DescribeTaskDefinitionInput, optFns ...func(*ecs.Options)) (*ecs.DescribeTaskDefinitionOutput, error) {
	return &m.Resp, nil
}

func TestDescribeTaskDefinition(t *testing.T) {
	resp := ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &ecstypes.TaskDefinition{
			Family: aws.String("dummy"),
		},
	}
	taskDefinition := &TaskDefinition{
		awsECS: mockedDescribeTaskDefinition{Resp: resp},
	}
	output, err := taskDefinition.DescribeTaskDefinition(context.Background(), "dummy")
	if err != nil {
		t.Error(err)
	}

	if *output.Family != "dummy" {
		t.Error("Task definition is invalid")
	}
}

func TestGetLogGroup(t *testing.T) {
	taskContainer := ecstypes.ContainerDefinition{
		Name: aws.String("TaskContainer"),
		LogConfiguration: &ecstypes.LogConfiguration{
			LogDriver: ecstypes.LogDriverAwslogs,
			Options: map[string]string{
				"awslogs-group":         "GroupName",
				"awslogs-stream-prefix": "LogPrefix",
			},
		},
	}
	dummyContainer := ecstypes.ContainerDefinition{
		Name: aws.String("DummyContainer"),
	}
	taskDef := &ecstypes.TaskDefinition{
		ContainerDefinitions: []ecstypes.ContainerDefinition{taskContainer, dummyContainer},
	}
	taskDefinition := &TaskDefinition{}
	group, prefix, err := taskDefinition.GetLogGroup(taskDef, "TaskContainer")
	if err != nil {
		t.Error(err)
	}

	if group != "GroupName" {
		t.Error("Group name is invalid")
	}
	if prefix != "LogPrefix" {
		t.Error("Stream prefix is invalid")
	}
}
