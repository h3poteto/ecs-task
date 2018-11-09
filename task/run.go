package task

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
