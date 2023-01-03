package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/littlehawk93/columba/config"
	"github.com/littlehawk93/columba/handler"
	"github.com/littlehawk93/columba/providers/fedex"
	"github.com/littlehawk93/columba/tracking"
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

	prov := &fedex.Provider{}

	events, err := prov.GetTrackingEvents("621648671050")

	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(&events, "", "  ")
	fmt.Println(string(b))

	return

	db, err := configuration.Database.Open()

	if err != nil {
		log.Fatalf("Unable to open database: %s\n", err.Error())
	}

	if err = tracking.Migrate(db); err != nil {
		log.Fatalf("Unable to set up database: %s\n", err.Error())
	}

	mainRouter := mux.NewRouter()

	handler.SetConfiguration(configuration)
	handler.AddAPIHandlers(mainRouter)

	fileServer := http.FileServer(http.Dir(configuration.WebRoot))
	mainRouter.PathPrefix("/").Handler(fileServer)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", configuration.Web.BindAddress, configuration.Web.Port), mainRouter); err != nil {
		log.Fatal(err)
	}
}
