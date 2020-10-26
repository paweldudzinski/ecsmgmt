package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Service struct gathers information about ECS service.
type Service struct {
	index          int
	arn            string
	name           string
	desiredCount   int64
	events         []*ecs.ServiceEvent
	pendingCount   int64
	runningCount   int64
	status         string
	taskDefinition string
}

// ListServices list ECS services within a cluster
func ListServices(session *session.Session, cluster string) ([]Service, error) {
	svc := ecs.New(session)
	input := &ecs.ListServicesInput{
		Cluster: aws.String(cluster),
	}
	result, err := svc.ListServices(input)
	if err != nil {
		return []Service{}, err
	}

	var services = []Service{}
	for i, arn := range result.ServiceArns {
		name := strings.Split(*arn, "/")[1]
		description := DescribeService(svc, name, cluster)
		s := Service{
			index:          i + 1,
			arn:            *arn,
			name:           name,
			desiredCount:   *description.DesiredCount,
			pendingCount:   *description.PendingCount,
			runningCount:   *description.RunningCount,
			status:         *description.Status,
			taskDefinition: strings.Split(*description.TaskDefinition, "/")[1],
		}
		if len(description.Events) > 20 {
			s.events = description.Events[:20]
		} else {
			s.events = description.Events
		}

		services = append(services, s)
	}
	return services, nil
}

// DescribeService returns service details
func DescribeService(svc *ecs.ECS, serviceName string, clusterName string) *ecs.Service {
	input := &ecs.DescribeServicesInput{
		Services: []*string{
			aws.String(serviceName),
		},
		Cluster: aws.String(clusterName),
	}

	result, _ := svc.DescribeServices(input)
	return result.Services[0]
}
