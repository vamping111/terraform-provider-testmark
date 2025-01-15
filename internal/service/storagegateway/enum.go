package storagegateway

import "time"

const (
	authenticationActiveDirectory = "ActiveDirectory"
	authenticationGuestAccess     = "GuestAccess"
)

func authentication_Values() []string {
	return []string{
		authenticationActiveDirectory,
		authenticationGuestAccess,
	}
}

const (
	bandwidthTypeAll      = "ALL"
	bandwidthTypeDownload = "DOWNLOAD"
	bandwidthTypeUpload   = "UPLOAD"
)

const (
	defaultStorageClassS3IntelligentTiering = "S3_INTELLIGENT_TIERING"
	defaultStorageClassS3OneZoneIA          = "S3_ONEZONE_IA"
	defaultStorageClassS3Standard           = "S3_STANDARD"
	defaultStorageClassS3StandardIA         = "S3_STANDARD_IA"
)

func defaultStorageClass_Values() []string {
	return []string{
		defaultStorageClassS3IntelligentTiering,
		defaultStorageClassS3OneZoneIA,
		defaultStorageClassS3Standard,
		defaultStorageClassS3StandardIA,
	}
}

const (
	gatewayTypeCached     = "CACHED"
	gatewayTypeFileFSXSMB = "FILE_FSX_SMB"
	gatewayTypeFileS3     = "FILE_S3"
	gatewayTypeStored     = "STORED"
	gatewayTypeVTL        = "VTL"
	gatewayTypeVTLSnow    = "VTL_SNOW"
)

func gatewayType_Values() []string {
	return []string{
		gatewayTypeCached,
		gatewayTypeFileFSXSMB,
		gatewayTypeFileS3,
		gatewayTypeStored,
		gatewayTypeVTL,
		gatewayTypeVTLSnow,
	}
}

const (
	mediumChangerTypeAWS_Gateway_VTL   = "AWS-Gateway-VTL"
	mediumChangerTypeIBM_03584L32_0402 = "IBM-03584L32-0402"
	mediumChangerTypeSTK_L700          = "STK-L700"
)

func mediumChangerType_Values() []string {
	return []string{
		mediumChangerTypeAWS_Gateway_VTL,
		mediumChangerTypeIBM_03584L32_0402,
		mediumChangerTypeSTK_L700,
	}
}

const (
	squashAllSquash  = "AllSquash"
	squashNoSquash   = "NoSquash"
	squashRootSquash = "RootSquash"
)

func squash_Values() []string {
	return []string{
		squashAllSquash,
		squashNoSquash,
		squashRootSquash,
	}
}

const (
	tapeDriveTypeIBM_ULT3580_TD5 = "IBM-ULT3580-TD5"
)

func tapeDriveType_Values() []string {
	return []string{
		tapeDriveTypeIBM_ULT3580_TD5,
	}
}

const (
	fileShareStatusAvailable     = "AVAILABLE"
	fileShareStatusCreating      = "CREATING"
	fileShareStatusDeleting      = "DELETING"
	fileShareStatusForceDeleting = "FORCE_DELETING"
	fileShareStatusUpdating      = "UPDATING"
)

const (
	fileSystemAssociationCreateTimeout = 10 * time.Minute
	fileSystemAssociationUpdateTimeout = 10 * time.Minute
	fileSystemAssociationDeleteTimeout = 10 * time.Minute
)

const (
	fileSystemAssociationStatusAvailable     = "AVAILABLE"
	fileSystemAssociationStatusCreating      = "CREATING"
	fileSystemAssociationStatusDeleting      = "DELETING"
	fileSystemAssociationStatusForceDeleting = "FORCE_DELETING"
	fileSystemAssociationStatusUpdating      = "UPDATING"
	fileSystemAssociationStatusError         = "ERROR"
)
