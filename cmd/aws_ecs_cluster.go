package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// Cluster struct gathers information about ECS cluster.
type Cluster struct {
	index                   int
	arn                     string
	name                    string
	status                  string // ACTIVE, PROVISIONING, DEPROVISIONING, FAILED, INACTIVE
	activeServicesCount     int64
	runningTasksCount       int64
	pendingTasksCount       int64
	containerInstancesCount int64
}

// ListClusters lists ECS clusters
func ListClusters(session *session.Session) ([]Cluster, error) {
	svc := ecs.New(session)
	input := &ecs.ListClustersInput{}
	result, err := svc.ListClusters(input)
	if err != nil {
		return []Cluster{}, err
	}

	var clusters []Cluster
	for i, arn := range result.ClusterArns {
		name := strings.Split(*arn, "/")[1]
		description := DescribeCluster(svc, name)
		cl := Cluster{
			index:                   i + 1,
			arn:                     *arn,
			name:                    name,
			status:                  *description.Status,
			activeServicesCount:     *description.ActiveServicesCount,
			runningTasksCount:       *description.RunningTasksCount,
			pendingTasksCount:       *description.PendingTasksCount,
			containerInstancesCount: *description.RegisteredContainerInstancesCount,
		}
		clusters = append(clusters, cl)
	}
	return clusters, nil
}

// DescribeCluster returns cluster details
func DescribeCluster(svc *ecs.ECS, clusterName string) *ecs.Cluster {
	input := &ecs.DescribeClustersInput{
		Clusters: []*string{
			aws.String(clusterName),
		},
	}
	result, _ := svc.DescribeClusters(input)
	return result.Clusters[0]
}
