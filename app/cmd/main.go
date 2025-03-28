package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"strings"
	"userservice/cmd/flags"
)

const (
	name = "user-microservice"

	defaultAllAddresses = "0.0.0.0"
)

var (
	commitHash string
	buildTime  string
	tag        string
)

var coreCmd = &cobra.Command{
	Use:     "user-microservice",
	Long:    "User microservice app binary",
	PreRunE: checkRequiredFlags,
	Run:     start,
}

func main() {
	configureCmd(coreCmd)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("GOMicroUserService") // TODO: review this naming

	//coreCmd.AddCommand(upgrader.Command())

	if err := coreCmd.Execute(); err != nil {
		// TODO: log error
		os.Exit(1)
	}
}

func start(cmd *cobra.Command, args []string) {}

func configureCmd(srcCmd *cobra.Command) {
	srcCmd.Flags().String(flags.ServerAddressFlag, defaultAllAddresses, "Listening address")
	err := viper.BindPFlag(flags.ServerAddressFlag, srcCmd.Flags().Lookup(flags.ServerAddressFlag))
	if err != nil {
		panic(err)
	}

	//srcCmd.Flags().Int(flags.RestAPIPortFlag, defaultHTTPPort, "Listening HTTP port")
	//err = viper.BindPFlag(flags.RestAPIPortFlag, srcCmd.Flags().Lookup(flags.RestAPIPortFlag))
	//if err != nil {
	//	panic(err)
	//}

	srcCmd.PersistentFlags().String(flags.ConfigPathFlag, ".", "Path to configuration file")
	err = viper.BindPFlag(flags.ConfigPathFlag, srcCmd.PersistentFlags().Lookup(flags.ConfigPathFlag))
	if err != nil {
		panic(err)
	}
}

func checkRequiredFlags(cmd *cobra.Command, _ []string) error {
	var requiredFlagsNotFound []string
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		required := len(flag.Annotations[cobra.BashCompOneRequiredFlag]) > 0 && flag.Annotations[cobra.BashCompOneRequiredFlag][0] == "true"
		if required && !envVarIsSet(flag.Name) {
			requiredFlagsNotFound = append(requiredFlagsNotFound, flag.Name)
		}
	})
	if len(requiredFlagsNotFound) > 0 {
		return fmt.Errorf("the following flags are required and were not present: %s ", strings.Join(requiredFlagsNotFound, ", "))
	}
	return nil
}
