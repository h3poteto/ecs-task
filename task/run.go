package task

import (
	"context"
	"regexp"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/service/ecs"
	log "github.com/sirupsen/logrus"
)

// Run a command on AWS ECS and output the log.
func (t *Task) Run() error {
	taskDef, err := t.taskDefinition.DescribeTaskDefinition(t.TaskDefinitionName)
	if err != nil {
		return err
	}
	group, streamPrefix, err := t.taskDefinition.GetLogGroup(taskDef, t.Container)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), t.Timeout)
	defer cancel()

	tasks, err := t.RunTask(ctx, taskDef)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	for _, task := range tasks {
		taskID := t.buildLogStream(task)
		w := NewWatcher(group, streamPrefix+"/"+t.Container+"/"+taskID, t.profile, t.region)
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := w.Polling(ctx)
			log.Warn(err)
		}()
	}

	err = t.WaitTask(ctx, tasks)
	time.Sleep(10 * time.Second)
	cancel()
	return nil
}

// buildLogStream returns a CloudWatchLog Stream name from ECS task.
// Task ARN format is `arn:aws:ecs:<region>:<aws_account_id>:task/c5cba4eb-5dad-405e-96db-71ef8eefe6a8`.
// And Log Stream format is `stream_prefix/container_name/task_id`.
func (t *Task) buildLogStream(task *ecs.Task) string {
	arn := *task.TaskArn
	taskRegexp := regexp.MustCompile(`task\/([a-z\d\-]+)`)
	return taskRegexp.FindStringSubmatch(arn)[1]
}
