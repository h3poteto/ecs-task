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
	ctx, cancel := context.WithCancel(context.Background())
	if t.Timeout != 0 {
		ctx, cancel = context.WithTimeout(ctx, t.Timeout)
	}
	defer cancel()

	task, err := t.RunTask(ctx, taskDef)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup

	taskID := t.buildLogStream(task)
	w := NewWatcher(group, streamPrefix+"/"+t.Container+"/"+taskID, t.profile, t.region, t.timestampFormat)
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := w.Polling(ctx)
		log.Error(err)
	}()

	err = t.WaitTask(ctx, task)
	

	time.Sleep(10 * time.Second)
	cancel()
	return err
}

// buildLogStream returns a CloudWatchLog Stream name from ECS task.
// Task ARN format is `arn:aws:ecs:<region>:<aws_account_id>:task(/<cluster_name>)/c5cba4eb-5dad-405e-96db-71ef8eefe6a8`.
// And Log Stream format is `stream_prefix/container_name/task_id`.
func (t *Task) buildLogStream(task *ecs.Task) string {
	arn := *task.TaskArn
	taskRegexp := regexp.MustCompile(`\/([a-z\d\-]+)$`)
	return taskRegexp.FindStringSubmatch(arn)[1]
}
