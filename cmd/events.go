package cmd

import (
	"github.com/spf13/cobra"
)

// eventsCmd represents the events command
var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Lists recent service events",
	Long:  `Lists recent 20 service events. `,
	Run: func(cmd *cobra.Command, args []string) {
		cluster, _ := cmd.Flags().GetString("cluster")
		service, _ := cmd.Flags().GetString("service")
		doListEvents(cmd, args, cluster, service)
	},
}

func init() {
	listCmd.AddCommand(eventsCmd)
	eventsCmd.Flags().String("cluster", "", "cluster name (required)")
	eventsCmd.Flags().String("service", "", "service name (required)")
	eventsCmd.MarkFlagRequired("cluster")
	eventsCmd.MarkFlagRequired("service")
}

func doListEvents(cmd *cobra.Command, args []string, cluster string, service string) {
	go Spinner()
	c := AWSCredentials{}
	c.InitAWS(cmd)
	session, _ := c.GetSession()
	services, err := ListServices(session, cluster)
	if err != nil {
		return
	}
	for _, s := range services {
		if s.name == service {
			PrintEvents(s.events)
		}
	}
}
