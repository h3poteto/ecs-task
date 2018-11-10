package task

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type mockedDescribeTaskDefinition struct {
	ecsiface.ECSAPI
	Resp ecs.DescribeTaskDefinitionOutput
}

func (m mockedDescribeTaskDefinition) DescribeTaskDefinition(in *ecs.DescribeTaskDefinitionInput) (*ecs.DescribeTaskDefinitionOutput, error) {
	return &m.Resp, nil
}

func TestDescribeTaskDefinition(t *testing.T) {
	resp := ecs.DescribeTaskDefinitionOutput{
		TaskDefinition: &ecs.TaskDefinition{
			Family: aws.String("dummy"),
		},
	}
	taskDefinition := &TaskDefinition{
		awsECS: mockedDescribeTaskDefinition{Resp: resp},
	}
	output, err := taskDefinition.DescribeTaskDefinition("dummy")
	if err != nil {
		t.Error(err)
	}

	if *output.Family != "dummy" {
		t.Error("Task definition is invalid")
	}
}

func TestGetLogGroup(t *testing.T) {
	taskContainer := &ecs.ContainerDefinition{
		Name: aws.String("TaskContainer"),
		LogConfiguration: &ecs.LogConfiguration{
			LogDriver: aws.String("awslogs"),
			Options: map[string]*string{
				"awslogs-group":         aws.String("GroupName"),
				"awslogs-stream-prefix": aws.String("LogPrefix"),
			},
		},
	}
	dummyContainer := &ecs.ContainerDefinition{
		Name: aws.String("DummyContainer"),
	}
	taskDef := &ecs.TaskDefinition{
		ContainerDefinitions: []*ecs.ContainerDefinition{taskContainer, dummyContainer},
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
