package cmd

import (
	"time"

	"github.com/h3poteto/ecs-task/task"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type runTask struct {
	cluster        string
	container      string
	taskDefinition string
	command        string
	timeout        int
}

func runTaskCmd() *cobra.Command {
	r := &runTask{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a task on ECS",
		Run:   r.run,
	}

	flags := cmd.Flags()
	flags.StringVarP(&r.cluster, "cluster", "c", "", "Name of ECS Cluster")
	flags.StringVar(&r.container, "container", "", "Name of container name in task definition")
	flags.StringVarP(&r.taskDefinition, "task-definition", "d", "", "Name of task definition to run task. Family and revision (family:revision), only Family or full ARN")
	flags.StringVar(&r.command, "command", "", "Command which you want to run")
	flags.IntVarP(&r.timeout, "timeout", "t", 0, "Timeout seconds")

	return cmd
}

func (r *runTask) run(cmd *cobra.Command, args []string) {
	profile, region, verbose := generalConfig()
	if !verbose {
		log.SetLevel(log.WarnLevel)
	}
	t, err := task.NewTask(r.cluster, r.container, r.taskDefinition, r.command, (time.Duration(r.timeout) * time.Second), profile, region)
	if err != nil {
		log.Fatal(err)
	}
	if err := t.Run(); err != nil {
		log.Fatal(err)
	}
}
