package cmd

import "github.com/spf13/cobra"

type runTask struct {
	cluster        string
	container      string
	taskDefinition string
	command        string
	timeout        int
}

func runTaskCmd() *cobra.Command {
	t := &runTask{}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a task on ECS",
		Run:   t.run,
	}

	flags := cmd.Flags()
	flags.StringVarP(&t.cluster, "cluster", "c", "", "Name of ECS Cluster")
	flags.StringVar(&t.container, "container", "", "Name of container name in task definition")
	flags.StringVarP(&t.taskDefinition, "task-definition", "d", "", "Name of task definition to run task. Family and revision (family:revision), only Family or full ARN")
	flags.StringVar(&t.command, "command", "", "Command which you want to run")
	flags.IntVarP(&t.timeout, "timeout", "t", 600, "Timeout seconds")

	return cmd
}

func (t *runTask) run(cmd *cobra.Command, args []string) {
}
