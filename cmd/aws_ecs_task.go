package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Task struct gathers information about task
type Task struct {
	arn                  string
	name                 string
	family               string
	status               string
	containerDefinitions []*ecs.ContainerDefinition
}

// ListTasks lists tasks definitions
func ListTasks(session *session.Session, family string) ([]Task, error) {
	input := &ecs.ListTaskDefinitionsInput{}
	svc := ecs.New(session)
	if family != "" {
		input = &ecs.ListTaskDefinitionsInput{
			FamilyPrefix: aws.String(family),
		}
	}

	result, err := svc.ListTaskDefinitions(input)
	if err != nil {
		return []Task{}, err
	}

	var tasks = []Task{}
	for _, arn := range result.TaskDefinitionArns {
		name := strings.Split(*arn, "/")[1]
		description := DescribeTask(svc, name)
		t := Task{
			name:                 name,
			arn:                  *arn,
			family:               *description.Family,
			status:               *description.Status,
			containerDefinitions: description.ContainerDefinitions,
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// DescribeTask returns task definition details
func DescribeTask(svc *ecs.ECS, name string) *ecs.TaskDefinition {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(name),
	}
	result, _ := svc.DescribeTaskDefinition(input)
	return result.TaskDefinition
}
