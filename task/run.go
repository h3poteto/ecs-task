package task

// Run a command on AWS ECS and output the log.
func (t *Task) Run() error {
	taskDef, err := t.taskDefinition.DescribeTaskDefinition(t.TaskDefinitionName)
	if err != nil {
		return err
	}
	_, err = t.RunTask(taskDef)
	if err != nil {
		return err
	}
	return nil
}
