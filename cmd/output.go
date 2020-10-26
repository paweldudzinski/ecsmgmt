package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/olekukonko/tablewriter"
)

var (
	captionDefaultColorScheme = []int{tablewriter.Normal, tablewriter.FgCyanColor}
	statusActiveColorScheme   = []int{tablewriter.Normal, tablewriter.FgGreenColor}
	statusInactiveColorScheme = []int{tablewriter.Normal, tablewriter.FgRedColor}
	headerClusterTitles       = []string{
		"",
		"Cluster Name",
		"Status",
		"Services Count",
		"Running Tasks",
		"Pending Tasks",
		"Container Instances",
	}
	headerServiceTitles = []string{
		"",
		"Service Name",
		"Status",
		"Task Definition",
		"Desired",
		"Running",
		"Pending",
	}
	headerEventsTitle = []string{
		"Date/time (UTC)",
		"Event",
	}
	headerInstancesTitle = []string{
		"",
		"Instance ID",
		"Status",
		"EC2 Public IP",
		"When registered (UTC)",
		"Running Tasks",
		"CPU",
		"Memory",
		"CPU Util %",
	}
	headerTaskTitle = []string{
		"Task Name",
		"Family",
		"Status",
		"Container name (Cpu/Memory)",
	}
)

const (
	statusIndex  = 2
	statusActive = "ACTIVE"
)

// PrintClusterInfo pretty prints clusters informantion
func PrintClusterInfo(c []Cluster) {
	data := [][]string{}
	for _, cluster := range c {
		row := []string{
			strconv.Itoa(cluster.index),
			cluster.name,
			cluster.status,
			strconv.Itoa(int(cluster.activeServicesCount)),
			strconv.Itoa(int(cluster.runningTasksCount)),
			strconv.Itoa(int(cluster.pendingTasksCount)),
			strconv.Itoa(int(cluster.containerInstancesCount)),
		}
		data = append(data, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headerClusterTitles)

	for _, v := range data {
		statusColorScheme := statusActiveColorScheme
		if v[statusIndex] != statusActive {
			statusColorScheme = statusInactiveColorScheme
		}
		table.Rich(v, []tablewriter.Colors{
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			statusColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
		})
	}
	fmt.Printf("\r")
	table.Render()
}

// PrintServicesInfo pretty prints services informantion
func PrintServicesInfo(s []Service) {
	data := [][]string{}
	for _, service := range s {
		row := []string{
			strconv.Itoa(service.index),
			service.name,
			service.status,
			service.taskDefinition,
			strconv.Itoa(int(service.desiredCount)),
			strconv.Itoa(int(service.runningCount)),
			strconv.Itoa(int(service.pendingCount)),
		}
		data = append(data, row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headerServiceTitles)

	for _, v := range data {
		statusColorScheme := statusActiveColorScheme
		if v[statusIndex] != statusActive {
			statusColorScheme = statusInactiveColorScheme
		}
		table.Rich(v, []tablewriter.Colors{
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			statusColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
		})
	}
	fmt.Printf("\r")
	table.Render()
}

// PrintEvents pretty print service events
func PrintEvents(events []*ecs.ServiceEvent) {
	data := [][]string{}
	for _, e := range events {
		row := []string{
			e.CreatedAt.Format("2006-01-02 15:04:05"),
			printEventMessage(*e.Message),
		}
		data = append(data, row)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headerEventsTitle)
	table.SetAutoWrapText(false)
	table.SetReflowDuringAutoWrap(false)

	for _, v := range data {
		table.Rich(v, []tablewriter.Colors{
			captionDefaultColorScheme,
			captionDefaultColorScheme,
		})
	}
	fmt.Printf("\r")
	table.Render()
}

func printEventMessage(message string) string {
	if len(message) > 100 {
		return message[:100] + "..."
	}
	return message
}

// PrintInstancesInfo pretty prints container instance description
func PrintInstancesInfo(i []ContainerInstance) {
	data := [][]string{}
	for _, instance := range i {
		row := []string{
			strconv.Itoa(instance.index),
			fmt.Sprintf("%v (%v)", instance.ec2InstanceID, instance.ec2InstanceType),
			instance.status,
			instance.publicIP,
			instance.whenRegistered.Format("2006-01-02 15:04:05"),
			strconv.Itoa(int(instance.runningTasks)),
			fmt.Sprintf("%v / %v", strconv.Itoa(int(instance.registeredResource.CPU)), strconv.Itoa(int(instance.remainingResource.CPU))),
			fmt.Sprintf("%v / %v", strconv.Itoa(int(instance.registeredResource.memory)), strconv.Itoa(int(instance.remainingResource.memory))),
			strconv.FormatFloat(instance.CPUUtilization, 'g', 4, 64),
		}
		data = append(data, row)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headerInstancesTitle)
	table.SetAutoWrapText(false)
	table.SetReflowDuringAutoWrap(false)

	for _, v := range data {
		statusColorScheme := statusActiveColorScheme
		if v[statusIndex] != statusActive {
			statusColorScheme = statusInactiveColorScheme
		}
		table.Rich(v, []tablewriter.Colors{
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			statusColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			captionDefaultColorScheme,
		})
	}
	fmt.Printf("\r")
	table.Render()
	fmt.Println("* CPU and memory values: <registered>/<remaining>")
}

// PrintTasksInfo pretty print task information
func PrintTasksInfo(t []Task) {
	data := [][]string{}
	for _, task := range t {
		row := []string{
			task.name,
			task.family,
			task.status,
			parseContainerDefinition(task.containerDefinitions),
		}
		data = append(data, row)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headerTaskTitle)
	table.SetAutoWrapText(false)
	table.SetReflowDuringAutoWrap(false)

	for _, v := range data {
		statusColorScheme := statusActiveColorScheme
		if v[statusIndex] != statusActive {
			statusColorScheme = statusInactiveColorScheme
		}
		table.Rich(v, []tablewriter.Colors{
			captionDefaultColorScheme,
			captionDefaultColorScheme,
			statusColorScheme,
			captionDefaultColorScheme,
		})
	}
	fmt.Printf("\r")
	table.Render()
}

func parseContainerDefinition(cd []*ecs.ContainerDefinition) string {
	var buffer bytes.Buffer
	for i, c := range cd {
		if *c.Cpu == 0 {
			buffer.WriteString(fmt.Sprintf("%v (auto/%v)", *c.Name, *c.Memory))
		} else {
			buffer.WriteString(fmt.Sprintf("%v (%v/%v)", *c.Name, *c.Cpu, *c.Memory))
		}
		if i > 0 {
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}
