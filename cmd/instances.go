package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// instancesCmd represents the instances command
var instancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "Lists cluster container instances",
	Long:  `Lists cluster container instances for EC2 launch type`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, _ := cmd.Flags().GetString("cluster")
		doListInstances(cmd, args, cluster)
	},
}

func init() {
	listCmd.AddCommand(instancesCmd)
	instancesCmd.Flags().String("cluster", "", "cluster name (required)")
	instancesCmd.MarkFlagRequired("cluster")
}

func doListInstances(cmd *cobra.Command, args []string, cluster string) {
	go Spinner()
	c := AWSCredentials{}
	c.InitAWS(cmd)
	session, _ := c.GetSession()
	if instances, err := ListInstances(session, cluster); err == nil {
		PrintInstancesInfo(instances)
	} else {
		fmt.Println(err.Error())
	}
}
