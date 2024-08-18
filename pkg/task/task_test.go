package task

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type mockedRunTask struct {
	ecsiface.ECSAPI
	Run      ecs.RunTaskOutput
	Describe ecs.DescribeTasksOutput
}

type mockedWaitTask struct {
	ecsiface.ECSAPI
	Describe ecs.DescribeTasksOutput
}

func (m mockedRunTask) RunTaskWithContext(ctx aws.Context, in *ecs.RunTaskInput, opts ...request.Option) (*ecs.RunTaskOutput, error) {
	return &m.Run, nil
}

func (m mockedRunTask) DescribeTasksWithContext(ctx aws.Context, in *ecs.DescribeTasksInput, options ...request.Option) (*ecs.DescribeTasksOutput, error) {
	return &m.Describe, nil
}

func (m mockedWaitTask) DescribeTasksWithContext(ctx aws.Context, in *ecs.DescribeTasksInput, options ...request.Option) (*ecs.DescribeTasksOutput, error) {
	return &m.Describe, nil
}

func TestRunTask(t *testing.T) {
	runTask := ecs.RunTaskOutput{
		Tasks: []*ecs.Task{
			&ecs.Task{
				ClusterArn:        aws.String("dummy-cluster"),
				TaskDefinitionArn: aws.String("task-definition-arn"),
				Overrides: &ecs.TaskOverride{
					ContainerOverrides: []*ecs.ContainerOverride{
						&ecs.ContainerOverride{
							Command: []*string{
								aws.String("echo"),
							},
							Name: aws.String("dummy"),
						},
					},
				},
			},
		},
	}
	describe := ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{
			&ecs.Task{
				DesiredStatus: aws.String("STOPPED"),
				LastStatus:    aws.String("STOPPED"),
				Containers: []*ecs.Container{
					&ecs.Container{
						ExitCode: aws.Int64(0),
					},
				},
			},
		},
	}
	task := &Task{
		awsECS: mockedRunTask{Run: runTask, Describe: describe},
		Command: []*string{
			aws.String("echo"),
		},
		Timeout: 10 * time.Second,
	}
	ctx := context.Background()
	_, err := task.RunTask(ctx, &ecs.TaskDefinition{
		TaskDefinitionArn: aws.String("task-definition-arn"),
	})
	if err != nil {
		t.Error(err)
	}
}

func TestWaitTask(t *testing.T) {
	describe := ecs.DescribeTasksOutput{
		Tasks: []*ecs.Task{
			&ecs.Task{
				DesiredStatus: aws.String("STOPPED"),
				LastStatus:    aws.String("STOPPED"),
				Containers: []*ecs.Container{
					&ecs.Container{
						Name:     aws.String("target"),
						ExitCode: aws.Int64(0),
					},
					&ecs.Container{
						Name:     aws.String("sidecar"),
						ExitCode: aws.Int64(137),
					},
				},
			},
		},
	}
	task := &Task{
		awsECS: mockedWaitTask{Describe: describe},
		Command: []*string{
			aws.String("echo"),
		},
		Container: "target",
		Timeout:   10 * time.Second,
	}
	ctx := context.Background()
	ecstask := ecs.Task{
			TaskArn: aws.String("test-arn"),
		}
	err := task.WaitTask(ctx, &ecstask)
	if err != nil {
		t.Error(err)
	}
}
