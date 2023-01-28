package task

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/pkg/errors"
)

// TaskDefinition has client of aws-sdk-go.
type TaskDefinition struct {
	awsECS ecsiface.ECSAPI
}

// NewTaskDefinition returns a new TaskDefinition struct, and initialize aws ecs API client.
func NewTaskDefinition(awsECS ecsiface.ECSAPI) *TaskDefinition {
	return &TaskDefinition{
		awsECS,
	}
}

// DescribeTaskDefinition gets a task definition.
// The family for the latest ACTIVE revision, family and revision (family:revision)
// for a specific revision in the family, or full Amazon Resource Name (ARN)
// of the task definition to describe.
func (d *TaskDefinition) DescribeTaskDefinition(taskDefinitionName string) (*ecs.TaskDefinition, error) {
	params := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskDefinitionName),
	}
	resp, err := d.awsECS.DescribeTaskDefinition(params)
	if err != nil {
		return nil, err
	}

	return resp.TaskDefinition, nil
}

// GetLogGroup gets cloudwatch logs group and stream prefix.
func (d *TaskDefinition) GetLogGroup(taskDef *ecs.TaskDefinition, containerName string) (string, string, error) {
	var containerDefinition *ecs.ContainerDefinition

	for _, c := range taskDef.ContainerDefinitions {
		if *c.Name == containerName {
			containerDefinition = c
		}
	}
	if containerDefinition == nil {
		return "", "", errors.New("Cannot find container")
	}
	if *containerDefinition.LogConfiguration.LogDriver != "awslogs" {
		return "", "", errors.New("Log driver is not awslogs")
	}
	logDriver := containerDefinition.LogConfiguration.Options
	group := logDriver["awslogs-group"]
	streamPrefix := logDriver["awslogs-stream-prefix"]
	return *group, *streamPrefix, nil
}
