[![CircleCI](https://circleci.com/gh/h3poteto/ecs-task.svg?style=svg)](https://circleci.com/gh/h3poteto/ecs-task)
[![GitHub release](http://img.shields.io/github/release/h3poteto/ecs-task.svg?style=flat-square)](https://github.com/h3poteto/ecs-task/releases)
[![GoDoc](https://godoc.org/github.com/h3poteto/ecs-task/task?status.svg)](https://godoc.org/github.com/h3poteto/ecs-task/task)

# ecs-task

`ecs-task` is a command line tool to run a task on the ECS. The feature is

- Wait for completion of the task execution
- Get logs from CloudWatch Logs and output in stream

This is a command line tool, but you can use `task` as a package.
So when you write own task execition script for AWS ECS, you can embed `task` package in your golang source code and customize task recipe.
Please check [godoc](https://godoc.org/github.com/h3poteto/ecs-task/task).

## Install
Get binary from GitHub:

```
$ wget https://github.com/h3poteto/ecs-task/releases/download/v0.2.0/ecs-task_v0.2.0_linux_amd64.zip
$ unzip ecs-task_v0.2.0_linux_amd64.zip
$ ./ecs-task help
Run a task on ECS

Usage:
  ecs-task [command]

Available Commands:
  help        Help about any command
  run         Run a task on ECS
  version     Print the version number

Flags:
  -h, --help             help for ecs-task
      --profile string   AWS profile (detault is none, and use environment variables)
      --region string    AWS region (default is none, and use AWS_DEFAULT_REGION)
  -v, --verbose          Enable verbose mode

Use "ecs-task [command] --help" for more information about a command.
```

## Usage
Please provide a command to `--command`.

```
$ ./ecs-task run --cluster=base-default-prd --container=task --task-definition=fascia-web-prd-task --command="echo 'hoge'" --region=ap-northeast-1
[2018-11-10 19:13:15 +0900 JST] hoge
```

And if the command is failed on ECS, `ecs-task` exit with error.
```
$ ./ecs-task run --cluster=base-default-prd --container=task --task-definition=fascia-web-prd-task --command="hoge" --region=ap-northeast-1
[2018-11-10 18:29:24 +0900 JST] ./entrypoint.sh: exec: line 13: hoge: not found
FATA[0015] exit code: 127
exit status 1
$ echo $?
1
```

If you want to run the task as Fargate, please provide fargate flag and your subnet IDs.

```
$ ./ecs-task run --cluster=base-default-prd --container=task --task-definition=fascia-web-prd-task --command='echo "hoge"' --fargate=true --subnets='subnet-12easdb,subnet-34asbdf' --region=ap-northeast-1
```

## AWS IAM Policy
Below is a basic IAM Policy required for ecs-task.

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "AllowUserToECSTask",
      "Effect": "Allow",
      "Action": [
        "ecs:DescribeTaskDefinition",
        "ecs:RunTask",
        "ecs:DescribeTasks",
        "ecs:ListTasks",
        "logs:DescribeLogStreams",
        "logs:GetLogEvents",
        "iam:PassRole"
      ],
      "Resource": "*"
    }
  ]
}
```

## License
The package is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
