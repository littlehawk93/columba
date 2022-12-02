package cmd

import (
	"log"
	"os"

	"github.com/littlehawk93/columba/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string
var configuration config.ApplicationConfiguration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "columba",
	Short: "launch the columba web app",
	Run:   run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "path to the config file")

	rootCmd.MarkPersistentFlagFilename("config")
	rootCmd.MarkPersistentFlagRequired("config")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Unable to read config file '%s': %s", configFile, err.Error())
	}

	configuration = config.ApplicationConfiguration{}

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("Unable to parse config file '%s': %s", configFile, err.Error())
	}
}

func run(cmd *cobra.Command, args []string) {

}
