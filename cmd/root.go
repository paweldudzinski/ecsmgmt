package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecsmgmt",
	Short: "ECS management tool",
	Long: `ecsmgmt lets you to manage your ECS clusters.
Please provide AWS credentials as flags or environment variables (default to any AWS cli tools):
export AWS_SECRET_ACCESS_KEY=***
export AWS_ACCESS_KEY_ID=***
export AWS_DEFAULT_REGION=***`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("aws_secret_access_key", "k", viper.GetString("AWS_SECRET_ACCESS_KEY"), "Your AWS Secret Access Key")
	rootCmd.PersistentFlags().StringP("aws_access_key_id", "a", viper.GetString("AWS_ACCESS_KEY_ID"), "Your AWS Access Key ID")
	rootCmd.PersistentFlags().StringP("aws_default_region", "r", viper.GetString("AWS_DEFAULT_REGION"), "Your AWS Region")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ecsmgmt" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ecsmgmt")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
