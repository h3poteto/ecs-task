package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:           "ecs-task",
	Short:         "Run a task on ECS",
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().StringP("profile", "", "", "AWS profile (detault is none, and use environment variables)")
	RootCmd.PersistentFlags().StringP("region", "", "", "AWS region (default is none, and use AWS_DEFAULT_REGION)")
	viper.BindPFlag("profile", RootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("region", RootCmd.PersistentFlags().Lookup("region"))

	RootCmd.AddCommand(
		runTaskCmd(),
	)
}

func generalConfig() (string, string) {
	return viper.GetString("profile"), viper.GetString("region")
}
