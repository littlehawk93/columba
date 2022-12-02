package tracking

type PackageStatus uint8

const (
	PackageStatusActive   PackageStatus = 1
	PackageStatusArchived PackageStatus = 2
)
