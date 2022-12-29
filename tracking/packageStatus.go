package tracking

import "strings"

type PackageStatus int

const (
	PackageStatusError    PackageStatus = 0
	PackageStatusActive   PackageStatus = 1
	PackageStatusArchived PackageStatus = 2
)

// String returns this status as a string
func (me PackageStatus) String() string {

	switch me {
	case PackageStatusActive:
		return "Active"
	case PackageStatusArchived:
		return "Archived"
	}
	return "Error"
}

// ParsePackageStatus parses a string into a PackageStatus. Returns PackageStatusError if an unrecognized string is provided
func ParsePackageStatus(str string) PackageStatus {

	str = strings.TrimSpace(strings.ToLower(str))

	switch str {
	case "active":
		return PackageStatusActive
	case "archived":
		return PackageStatusArchived
	}
	return PackageStatusError
}
