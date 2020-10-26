package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// clustersCmd represents the clusters command
var clustersCmd = &cobra.Command{
	Use:   "clusters",
	Short: "Lists available clusters",
	Long:  `Lists available clusers with additional description`,
	Run: func(cmd *cobra.Command, args []string) {
		doListClusters(cmd, args)
	},
}

func init() {
	listCmd.AddCommand(clustersCmd)
}

// e.g. ecsmgmt list clusters
func doListClusters(cmd *cobra.Command, args []string) {
	go Spinner()
	c := AWSCredentials{}
	c.InitAWS(cmd)
	session, _ := c.GetSession()
	if clusters, err := ListClusters(session); err == nil {
		PrintClusterInfo(clusters)
	} else {
		fmt.Println(err.Error())
	}
}
