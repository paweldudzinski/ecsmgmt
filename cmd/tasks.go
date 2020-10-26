package cmd

import (
	"github.com/spf13/cobra"
)

// tasksCmd represents the tasks command
var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Lists tasks definitions",
	Long:  `Lists tasks definitions`,
	Run: func(cmd *cobra.Command, args []string) {
		family, _ := cmd.Flags().GetString("family")
		doListTasks(cmd, args, family)
	},
}

func init() {
	listCmd.AddCommand(tasksCmd)
	tasksCmd.Flags().String("family", "", "a task family prefix")
}

func doListTasks(cmd *cobra.Command, args []string, family string) {
	go Spinner()
	c := AWSCredentials{}
	c.InitAWS(cmd)
	session, _ := c.GetSession()
	if tasks, err := ListTasks(session, family); err == nil {
		PrintTasksInfo(tasks)
	}
}
