package providers

import (
	"regexp"
	"strings"

	"github.com/littlehawk93/columba/providers/fedex"
	"github.com/littlehawk93/columba/providers/ups"
	"github.com/littlehawk93/columba/providers/usps"
	"github.com/littlehawk93/columba/tracking"
)

var serviceProviders = map[string]tracking.Provider{
	"fedex": &fedex.Provider{},
	"usps":  &usps.Provider{},
	"ups":   &ups.Provider{},
}

// GetServiceProvider get a provider from a string, returns nil if the serviceID is invalid
func GetServiceProvider(serviceID string) tracking.Provider {

	serviceID = strings.ToLower(regexp.MustCompile(`[^a-zA-z]`).ReplaceAllString(serviceID, ""))

	if provider, ok := serviceProviders[serviceID]; !ok {
		return nil
	} else {
		return provider
	}
}

// GetAllServiceProviderNames get a slice of all the provider names supported by Columba
func GetAllServiceProviderNames() []string {

	names := make([]string, len(serviceProviders))

	for i := 0; i < len(names); i++ {
		names[i] = ""
	}

	for _, v := range serviceProviders {

		name := v.GetID()

		for i := 0; i < len(names); i++ {
			if names[i] == "" {
				names[i] = name
				break
			} else if names[i] > name {
				tmp := names[i]
				names[i] = name
				name = tmp
			}
		}
	}

	return names
}
