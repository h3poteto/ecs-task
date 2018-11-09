package task

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type TaskDefinition struct {
	awsECS ecsiface.ECSAPI
}

func NewTaskDefinition(profile, region string) *TaskDefinition {
	awsECS := ecs.New(session.New(), newConfig(profile, region))
	return &TaskDefinition{
		awsECS,
	}
}

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
