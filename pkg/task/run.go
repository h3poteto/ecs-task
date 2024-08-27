package task

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecstypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	log "github.com/sirupsen/logrus"
)

// Run a command on AWS ECS and output the log.
func (t *Task) Run() error {
	ctx := context.Background()
	taskDef, err := t.taskDefinition.DescribeTaskDefinition(ctx, t.TaskDefinitionName)
	if err != nil {
		return err
	}
	group, streamPrefix, err := t.taskDefinition.GetLogGroup(taskDef, t.Container)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	task, err := t.RunTask(ctx, taskDef)
	if err != nil {
		return err
	}

	taskID := t.buildLogStream(task)
	logPollDoneChan := make(chan struct{})
	pollLogsCtx, pollLogsCancel := context.WithCancel(ctx)
	w := NewWatcher(group, streamPrefix+"/"+t.Container+"/"+taskID, t.awsLogs, t.timestampFormat)
	go func() {
		defer close(logPollDoneChan)
		log.Info("Polling logs")
		err := w.Polling(pollLogsCtx)
		if err != nil {
			log.Errorf("Get logs thread failed: %v", err)
		} else {
			log.Info("Get logs thread gracefully stopping")
		}
	}()

	pollTaskStopDoneChan := make(chan error)
	pollExitCtx, pollExitCancel := context.WithCancel(ctx)
	defer pollExitCancel() // make go vet lostcancel happy
	go func() {
		defer close(pollTaskStopDoneChan)
		err := t.WaitTask(pollExitCtx, task)
		if err != nil {
			log.Errorf("Task status polling thread failed: %v", err)
		} else {
			log.Info("Task status polling thread gracefully stopping")
		}
		pollTaskStopDoneChan <- err
	}()

	var timeoutChan <-chan time.Time
	if t.Timeout > 0 {
		timeoutChan = time.After(t.Timeout)
	}

	var stopTaskReason string
	select {
	case sig := <-sigchan:
		log.WithFields(log.Fields{
			"signal": sig.String(),
		}).Info("Received signal; calling ecs.StopTask on task")
		stopTaskReason = fmt.Sprintf("ecs-task propagating signal %s", sig.String())
	case err = <-pollTaskStopDoneChan:
		log.Info("Task stopped on its own")
	case <-timeoutChan:
		log.WithFields(log.Fields{
			"timeout": t.Timeout,
		}).Info("Run timeout; calling ecs.StopTask on task")
	}
	if stopTaskReason != "" {
		params := ecs.StopTaskInput{
			Cluster: aws.String(t.Cluster),
			Reason:  aws.String(stopTaskReason),
			Task:    task.TaskArn,
		}
		if _, err := t.awsECS.StopTask(ctx, &params); err != nil {
			log.Errorf("Error calling ecs.StopTask: %v", err)
		}
		log.Info("After esc.StopTask; waiting up to 60s for task to stop")
		select {
		// wait for the default ECS_CONTAINER_STOP_TIMEOUT (=30s) + an additional 30s
		case <-time.After(60 * time.Second):
			log.Info("Task is still not done after 60s; giving up on checking its status")
			pollExitCancel()
			err = <-pollTaskStopDoneChan
		case err = <-pollTaskStopDoneChan:
		}
	}

	log.Info("Waiting 10s for more GetLogEvents")
	time.Sleep(10 * time.Second)
	log.Info("Shutting down get logs thread")
	pollLogsCancel()
	<-logPollDoneChan
	log.Info("Exiting")
	return err
}

// buildLogStream returns a CloudWatchLog Stream name from ECS task.
// Task ARN format is `arn:aws:ecs:<region>:<aws_account_id>:task(/<cluster_name>)/c5cba4eb-5dad-405e-96db-71ef8eefe6a8`.
// And Log Stream format is `stream_prefix/container_name/task_id`.
func (t *Task) buildLogStream(task *ecstypes.Task) string {
	arn := *task.TaskArn
	taskRegexp := regexp.MustCompile(`\/([a-z\d\-]+)$`)
	return taskRegexp.FindStringSubmatch(arn)[1]
}
