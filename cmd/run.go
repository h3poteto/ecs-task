package cmd

import (
	"time"

	"github.com/h3poteto/ecs-task/pkg/task"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type runTask struct {
	cluster         string
	container       string
	taskDefinition  string
	command         string
	subnets         string
	securityGroups  string
	fargate         bool
	timeout         int
	timestampFormat string
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
	flags.StringVarP(&r.subnets, "subnets", "s", "", "Provide subnet IDs with comma-separated string (subnet-12abcde,subnet-34abcde). This param is necessary, if you set farage flag.")
	flags.StringVarP(&r.securityGroups, "security-groups", "g", "", "Provide security group IDs with comma-separated string (sg-0123asdb,sg-2345asdf), if you want to attach the security groups to ENI of the task.")
	flags.BoolVarP(&r.fargate, "fargate", "f", false, "Whether run task with FARGATE")
	flags.IntVarP(&r.timeout, "timeout", "t", 0, "Timeout seconds")
	flags.StringVarP(&r.timestampFormat, "timestamp-format", "", "[2006-01-02 15:04:05.999999999 -0700 MST]", "Format of timestamp for outputs. You should follow the style of Time.Format (see https://golang.org/pkg/time/#pkg-constants)")

	return cmd
}

func (r *runTask) run(cmd *cobra.Command, args []string) {
	profile, region, verbose := generalConfig()
	if !verbose {
		log.SetLevel(log.WarnLevel)
	}
	t, err := task.NewTask(r.cluster, r.container, r.taskDefinition, r.command, r.fargate, r.subnets, r.securityGroups, (time.Duration(r.timeout) * time.Second), r.timestampFormat, profile, region)
	if err != nil {
		log.Fatal(err)
	}
	if err := t.Run(); err != nil {
		log.Fatal(err)
	}
}
