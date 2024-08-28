package task

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type mockedRunTask struct {
	ECSClient
	Run      ecs.RunTaskOutput
	Describe ecs.DescribeTasksOutput
}

type mockedWaitTask struct {
	ECSClient
	Describe ecs.DescribeTasksOutput
}

func (m mockedRunTask) RunTask(ctx context.Context, params *ecs.RunTaskInput, opts ...func(*ecs.Options)) (*ecs.RunTaskOutput, error) {
	return &m.Run, nil
}

func (m mockedRunTask) DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, options ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error) {
	return &m.Describe, nil
}

func (m mockedWaitTask) DescribeTasks(ctx context.Context, params *ecs.DescribeTasksInput, options ...func(*ecs.Options)) (*ecs.DescribeTasksOutput, error) {
	return &m.Describe, nil
}

func TestRunTask(t *testing.T) {
	runTask := ecs.RunTaskOutput{
		Tasks: []ecstypes.Task{
			ecstypes.Task{
				ClusterArn:        aws.String("dummy-cluster"),
				TaskDefinitionArn: aws.String("task-definition-arn"),
				Overrides: &ecstypes.TaskOverride{
					ContainerOverrides: []ecstypes.ContainerOverride{
						ecstypes.ContainerOverride{
							Command: []string{
								"echo",
							},
							Name: aws.String("dummy"),
						},
					},
				},
			},
		},
	}
	describe := ecs.DescribeTasksOutput{
		Tasks: []ecstypes.Task{
			ecstypes.Task{
				DesiredStatus: aws.String("STOPPED"),
				LastStatus:    aws.String("STOPPED"),
				Containers: []ecstypes.Container{
					ecstypes.Container{
						ExitCode: aws.Int32(0),
					},
				},
			},
		},
	}
	task := &Task{
		awsECS: mockedRunTask{Run: runTask, Describe: describe},
		Command: []string{
			"echo",
		},
		Timeout: 10 * time.Second,
	}
	ctx := context.Background()
	_, err := task.RunTask(ctx, &ecstypes.TaskDefinition{
		TaskDefinitionArn: aws.String("task-definition-arn"),
	})
	if err != nil {
		t.Error(err)
	}
}

func TestWaitTask(t *testing.T) {
	describe := ecs.DescribeTasksOutput{
		Tasks: []ecstypes.Task{
			ecstypes.Task{
				DesiredStatus: aws.String("STOPPED"),
				LastStatus:    aws.String("STOPPED"),
				Containers: []ecstypes.Container{
					ecstypes.Container{
						Name:     aws.String("target"),
						ExitCode: aws.Int32(0),
					},
					ecstypes.Container{
						Name:     aws.String("sidecar"),
						ExitCode: aws.Int32(137),
					},
				},
			},
		},
	}
	task := &Task{
		awsECS: mockedWaitTask{Describe: describe},
		Command: []string{
			"echo",
		},
		Container: "target",
		Timeout:   10 * time.Second,
	}
	ctx := context.Background()
	ecstask := ecstypes.Task{
		TaskArn: aws.String("test-arn"),
	}
	err := task.WaitTask(ctx, &ecstask)
	if err != nil {
		t.Error(err)
	}
}
