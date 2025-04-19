package bucket

type StorageTierEnum string

const (
	STierStandard    StorageTierEnum = "STANDARD"
	STierLowAccess   StorageTierEnum = "LOW"
	STierTierArchive StorageTierEnum = "ARCHIVE"
)
