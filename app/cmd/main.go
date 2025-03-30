package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"userservice/cmd/config"
	"userservice/cmd/flags"
	"userservice/internal/graph"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	name      = "usersvc"
	envPrefix = "SVC_USER"

	defaultServerAddress = "localhost"
	defaultServerPort    = 8080
	// This is coupled to this test-project, in a real app the default config path should make more sense
	defaultConfigPath = "../test/examples"
)

var rootCmd = &cobra.Command{
	Use:   name,
	Short: "go user microservice",
	Run:   startService,
}

func main() {
	configureCmd(rootCmd)
	configureEnvVariables()

	if err := rootCmd.Execute(); err != nil {
		// TODO: refactor this to a log entry
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func configureEnvVariables() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.SetEnvPrefix(envPrefix)
}

func configureCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String(flags.ServerAddressFlag, defaultServerAddress, "Server address")
	cmd.PersistentFlags().Int(flags.PortFlag, defaultServerPort, "Port to listen on")
	cmd.PersistentFlags().String(flags.ConfigPathFlag, defaultConfigPath, "Path to the service configuration file")

	err := viper.BindPFlag(flags.ServerAddressFlag, cmd.PersistentFlags().Lookup(flags.ServerAddressFlag))
	if err != nil {
		panic(fmt.Sprintf("Error binding ServerAddressFlag: %s\n", err))
	}
	err = viper.BindPFlag(flags.PortFlag, cmd.PersistentFlags().Lookup(flags.PortFlag))
	if err != nil {
		panic(fmt.Sprintf("Error binding PortFlag: %s\n", err))
	}
	err = viper.BindPFlag(flags.ConfigPathFlag, cmd.PersistentFlags().Lookup(flags.ConfigPathFlag))
	if err != nil {
		panic(fmt.Sprintf("Error binding ConfigPathFlag: %s\n", err))
	}

	cobra.OnInitialize(ensureConfigFileExists)
}

func ensureConfigFileExists() {
	configPath := viper.GetString(flags.ConfigPathFlag)
	viper.SetConfigType("yaml")

	_, err := os.Stat(configPath)
	if err != nil {
		panic(fmt.Sprintf("Config file [%s] does not exist\n", configPath))
	}

	viper.SetConfigFile(configPath)
}

func startService(_ *cobra.Command, _ []string) {
	address := viper.GetString(flags.ServerAddressFlag)
	port := viper.GetInt(flags.PortFlag)
	configPath := viper.GetString(flags.ConfigPathFlag)

	graphConfig, err := config.GetGraphConfig(configPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %s\n", err))
	}

	service, err := graph.New(*graphConfig)
	if err != nil {
		panic(fmt.Sprintf("Error creating service graph: %s\n", err))
	}

	// Build injects all the required dependencies
	service.Build()

	err = service.StartServer(graph.StartServerOptions{
		Port: strconv.Itoa(port),
		Host: address,
	})
	if err != nil {
		panic(fmt.Sprintf("Error starting service server: %s\n", err))
	}
}
