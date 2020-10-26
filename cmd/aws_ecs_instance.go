package cmd

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Resource contains information about CPU and memory resources
type Resource struct {
	CPU    int64
	memory int64
}

// ContainerInstance contains information about EC2 instances
type ContainerInstance struct {
	index              int
	uuid               string
	runningTasks       int64
	whenRegistered     time.Time
	status             string
	ec2InstanceID      string
	ec2InstanceType    string
	registeredResource Resource
	remainingResource  Resource
	publicIP           string
	CPUUtilization     float64
}

// ListInstances lists cluster container instances
func ListInstances(session *session.Session, cluster string) ([]ContainerInstance, error) {
	svc := ecs.New(session)
	input := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(cluster),
	}

	result, err := svc.ListContainerInstances(input)
	if err != nil {
		return []ContainerInstance{}, err
	}

	var containerInstances = []ContainerInstance{}
	for i, arn := range result.ContainerInstanceArns {
		description := DescribeContainerInstances(svc, *arn, cluster)
		ec2, _ := DescribeEC2Instance(session, *description.Ec2InstanceId)
		cpuAvg := GetEC2Metrics(session, *description.Ec2InstanceId)
		c := ContainerInstance{
			index:              i + 1,
			uuid:               *arn,
			whenRegistered:     *description.RegisteredAt,
			status:             *description.Status,
			ec2InstanceID:      *description.Ec2InstanceId,
			ec2InstanceType:    *ec2.InstanceType,
			registeredResource: getResources(description.RegisteredResources),
			remainingResource:  getResources(description.RemainingResources),
			runningTasks:       *description.RunningTasksCount,
			publicIP:           *ec2.PublicIpAddress,
			CPUUtilization:     cpuAvg,
		}
		containerInstances = append(containerInstances, c)
	}
	return containerInstances, nil
}

// DescribeContainerInstances returns description of a container instance by uuid
func DescribeContainerInstances(svc *ecs.ECS, uuid string, clusterName string) *ecs.ContainerInstance {
	input := &ecs.DescribeContainerInstancesInput{
		Cluster: aws.String(clusterName),
		ContainerInstances: []*string{
			aws.String(uuid),
		},
	}

	result, _ := svc.DescribeContainerInstances(input)
	return result.ContainerInstances[0]
}

func getResources(resources []*ecs.Resource) Resource {
	var (
		rCPU int64
		rMem int64
	)
	for _, r := range resources {
		if *r.Name == "CPU" {
			rCPU = *r.IntegerValue
		} else if *r.Name == "MEMORY" {
			rMem = *r.IntegerValue
		}
	}
	return Resource{
		CPU:    rCPU,
		memory: rMem,
	}
}

// DescribeEC2Instance desc
func DescribeEC2Instance(session *session.Session, id string) (*ec2.Instance, error) {
	svc := ec2.New(session)
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(id),
		},
	}

	result, _ := svc.DescribeInstances(input)
	r := result.Reservations[0]
	return r.Instances[0], nil
}

// GetEC2Metrics gets CloundWatch metrics for a given instance
func GetEC2Metrics(session *session.Session, id string) float64 {
	svc := cloudwatch.New(session)
	period := 360 * time.Second
	now := time.Now().Local().UTC()
	end := now.Add(-period)
	start := end.Add(-period)

	input := &cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{
			&cloudwatch.Dimension{
				Name:  aws.String("InstanceId"),
				Value: aws.String(id),
			},
		},
		MetricName: aws.String("CPUUtilization"),
		Namespace:  aws.String("AWS/EC2"),
		StartTime:  aws.Time(start),
		EndTime:    aws.Time(end),
		Statistics: []*string{
			aws.String(cloudwatch.StatisticAverage),
		},
		Period: aws.Int64(int64(period.Seconds())),
	}
	stats, _ := svc.GetMetricStatistics(input)
	if points := stats.Datapoints; len(points) != 0 {
		return *points[0].Average
	}
	return 0.0
}
