package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/littlehawk93/columba/providers"
	"github.com/littlehawk93/columba/providers/utils"
	"github.com/spf13/cobra"
)

var testProviderName string
var testTrackingNumber string

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test a tracking provider and number and view events parsed",
	Run:   executeTest,
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().StringVarP(&testTrackingNumber, "tracknum", "n", "", "The tracking number to test with")
	testCmd.Flags().StringVarP(&testProviderName, "provider", "p", "", "The service provider to get tracking data from")

	testCmd.MarkFlagRequired("tracknum")
	testCmd.MarkFlagRequired("provider")
}

func executeTest(cmd *cobra.Command, args []string) {

	utils.DebugRequests = true
	prov := providers.GetServiceProvider(testProviderName)

	if prov == nil {
		log.Fatalf("Invalid provider name '%s'\n", testProviderName)
	}

	events, err := prov.GetTrackingEvents(testTrackingNumber)

	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(&events, "", "  ")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))
}
