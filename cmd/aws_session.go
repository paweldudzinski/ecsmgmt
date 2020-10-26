package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AWSCredentials is a collection of AWS credentials
type AWSCredentials struct {
	SecretAccessKey string
	AccessKeyID     string
	DefaultRegion   string
}

// InitAWS collects AWS credentials from cmd line or env variables.
func (c *AWSCredentials) InitAWS(cmd *cobra.Command) {
	awsSecretAccessKeyEnv := viper.GetString("aws_secret_access_key")
	awsSecretAccessKeyCmd, _ := cmd.Flags().GetString("aws_secret_access_key")
	if awsSecretAccessKeyEnv != "" {
		c.SecretAccessKey = awsSecretAccessKeyEnv
	}
	if awsSecretAccessKeyCmd != "" {
		c.SecretAccessKey = awsSecretAccessKeyCmd
	}

	awsAccessKeyIDEnv := viper.GetString("aws_access_key_id")
	awsAccessKeyIDCmd, _ := cmd.Flags().GetString("aws_access_key_id")
	if awsAccessKeyIDEnv != "" {
		c.AccessKeyID = awsAccessKeyIDEnv
	}
	if awsAccessKeyIDCmd != "" {
		c.AccessKeyID = awsAccessKeyIDCmd
	}

	awsDefaultRegionEnv := viper.GetString("aws_default_region")
	awsDefaultRegionCmd, _ := cmd.Flags().GetString("aws_default_region")
	if awsDefaultRegionEnv != "" {
		c.DefaultRegion = awsDefaultRegionEnv
	}
	if awsDefaultRegionCmd != "" {
		c.DefaultRegion = awsDefaultRegionCmd
	}
}

// GetSession returns AWS session identifier
func (c *AWSCredentials) GetSession() (*session.Session, error) {
	session, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessKeyID, c.SecretAccessKey, ""),
		Region:      aws.String(c.DefaultRegion),
	})
	return session, err
}
