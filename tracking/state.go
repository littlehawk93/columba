package tracking

import "strings"

var statesList []string = []string{
	"AL",
	"AK",
	"AZ",
	"AR",
	"CA",
	"CO",
	"CT",
	"DE",
	"FL",
	"GA",
	"HI",
	"ID",
	"IL",
	"IN",
	"IA",
	"KS",
	"KY",
	"LA",
	"ME",
	"MD",
	"MA",
	"MI",
	"MN",
	"MS",
	"MO",
	"MT",
	"NE",
	"NV",
	"NH",
	"NJ",
	"NM",
	"NY",
	"NC",
	"ND",
	"OH",
	"OK",
	"OR",
	"PA",
	"RI",
	"SC",
	"SD",
	"TN",
	"TX",
	"UT",
	"VT",
	"VA",
	"WA",
	"WV",
	"WI",
	"WY",
}

// GetAllStates returns a set of all states in abbreviated form
func GetAllStates() []string {
	states := make([]string, len(statesList))
	copy(states, statesList)
	return states
}

// IsStateAbbreviation returns true if the provided string is an abbreviated state name (case insensitive), false otherwise
func IsStateAbbreviation(state string) bool {

	state = strings.ToUpper(strings.TrimSpace(state))

	if len(state) != 2 {
		return false
	}

	for _, s := range statesList {
		if state == s {
			return true
		}
	}
	return false
}
