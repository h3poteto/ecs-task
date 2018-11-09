package task

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Task struct {
	awsECS ecsiface.ECSAPI

	Cluster            string
	Container          string
	TaskDefinitionName string
	taskDefinition     *TaskDefinition
	Command            []*string
	Timeout            time.Duration
}

func NewTask(cluster, container, taskDefinitionName, command string, timeout time.Duration, profile, region string) (*Task, error) {
	if cluster == "" {
		return nil, errors.New("Cluster name is required")
	}
	if container == "" {
		return nil, errors.New("Container name is required")
	}
	if taskDefinitionName == "" {
		return nil, errors.New("Task definition is required")
	}
	if command == "" {
		return nil, errors.New("Comamnd is reqired")
	}
	awsECS := ecs.New(session.New(), newConfig(profile, region))
	taskDefinition := NewTaskDefinition(profile, region)
	p := shellwords.NewParser()
	commands, err := p.Parse(command)
	if err != nil {
		return nil, errors.Wrap(err, "Parse error")
	}
	var cmd []*string
	for _, c := range commands {
		cmd = append(cmd, aws.String(c))
	}

	return &Task{
		awsECS:             awsECS,
		Cluster:            cluster,
		Container:          container,
		TaskDefinitionName: taskDefinitionName,
		taskDefinition:     taskDefinition,
		Command:            cmd,
		Timeout:            timeout,
	}, nil
}

// RunTask calls run-task API.
func (t *Task) RunTask(taskDefinition *ecs.TaskDefinition) ([]*ecs.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()

	containerOverride := &ecs.ContainerOverride{
		Command: t.Command,
		Name:    aws.String(t.Container),
	}

	override := &ecs.TaskOverride{
		ContainerOverrides: []*ecs.ContainerOverride{
			containerOverride,
		},
	}

	params := &ecs.RunTaskInput{
		Cluster:        aws.String(t.Cluster),
		TaskDefinition: taskDefinition.TaskDefinitionArn,
		Overrides:      override,
	}
	resp, err := t.awsECS.RunTaskWithContext(ctx, params)
	if err != nil {
		return nil, err
	}
	if len(resp.Failures) > 0 {
		log.Errorf("Run task error: %+v", resp.Failures)
		return nil, errors.New(*resp.Failures[0].Reason)
	}
	log.Infof("Running tasks: %+v", resp.Tasks)

	err = t.waitRunning(ctx, resp.Tasks)
	if err != nil {
		return resp.Tasks, err
	}
	return resp.Tasks, nil
}

// waitRunning waits a task running.
func (t *Task) waitRunning(ctx context.Context, tasks []*ecs.Task) error {
	log.Info("Waiting for running task...")

	taskArns := []*string{}
	for _, task := range tasks {
		taskArns = append(taskArns, task.TaskArn)
	}
	errCh := make(chan error, 1)
	done := make(chan struct{}, 1)
	go func() {
		err := t.waitExitTasks(taskArns)
		if err != nil {
			errCh <- err
		}
		close(done)
	}()
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-done:
		log.Info("Run task is success")
	case <-ctx.Done():
		return errors.New("process timeout")
	}

	return nil
}

func (t *Task) waitExitTasks(taskArns []*string) error {
retry:
	for {
		time.Sleep(5 * time.Second)

		params := &ecs.DescribeTasksInput{
			Cluster: aws.String(t.Cluster),
			Tasks:   taskArns,
		}
		resp, err := t.awsECS.DescribeTasks(params)
		if err != nil {
			return err
		}

		for _, task := range resp.Tasks {
			if !t.checkTaskStopped(task) {
				continue retry
			}
		}

		for _, task := range resp.Tasks {
			code, result, err := t.checkTaskSucceeded(task)
			if err != nil {
				continue retry
			}
			if !result {
				return errors.Errorf("exit code: %v", code)
			}
		}
		return nil
	}
}

func (t *Task) checkTaskStopped(task *ecs.Task) bool {
	if *task.DesiredStatus != "STOPPED" {
		return false
	}
	return true
}

func (t *Task) checkTaskSucceeded(task *ecs.Task) (int64, bool, error) {
	for _, c := range task.Containers {
		if c.ExitCode == nil {
			return 1, false, errors.New("can not read exit code")
		}
		if *c.ExitCode != int64(0) {
			return *c.ExitCode, false, nil
		}
	}
	return int64(0), true, nil
}
