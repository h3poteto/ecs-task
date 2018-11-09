package task

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

// TaskDefinition has client of aws-sdk-go.
type TaskDefinition struct {
	awsECS ecsiface.ECSAPI
}

// NewTaskDefinition returns a new TaskDefinition struct, and initialize aws ecs API client.
func NewTaskDefinition(profile, region string) *TaskDefinition {
	awsECS := ecs.New(session.New(), newConfig(profile, region))
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
