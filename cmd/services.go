package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// servicesCmd represents the services command
var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Lists cluster's active services",
	Long:  `Lists cluster's active services with details`,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, _ := cmd.Flags().GetString("cluster")
		doSListervices(cmd, args, cluster)
	},
}

func init() {
	listCmd.AddCommand(servicesCmd)
	servicesCmd.Flags().String("cluster", "", "cluster name (required)")
	servicesCmd.MarkFlagRequired("cluster")
}

// e.g. ecsctl get service --cluster <cluster_name>
func doSListervices(cmd *cobra.Command, args []string, cluster string) {
	go Spinner()
	c := AWSCredentials{}
	c.InitAWS(cmd)
	session, _ := c.GetSession()
	if services, err := ListServices(session, cluster); err == nil {
		PrintServicesInfo(services)
	} else {
		fmt.Println(err.Error())
	}
}
