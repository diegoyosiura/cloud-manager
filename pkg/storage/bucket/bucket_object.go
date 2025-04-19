package bucket

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"time"
)

type BucketObject struct {
	Key          string
	LastModified time.Time
	Size         int64
	StorageClass StorageTierEnum
}

func NewBucketObjectFromAWS(o *s3.Object) BucketObject {
	var tier StorageTierEnum
	switch *o.StorageClass {
	case "STANDARD":
	case "STANDARD_IA":
	case "EXPRESS_ONEZONE":
		tier = STierStandard
		break
	case "REDUCED_REDUNDANCY":
	case "INTELLIGENT_TIERING":
		tier = STierLowAccess
		break
	case "DEEP_ARCHIVE":
	case "GLACIER":
	case "GLACIER_IR":
	case "ONEZONE_IA":
		tier = STierTierArchive
		break
	}

	return BucketObject{
		Key:          *o.Key,
		LastModified: *o.LastModified,
		Size:         *o.Size,
		StorageClass: tier,
	}
}

func NewBucketObjectFromOCI(o objectstorage.ObjectSummary) BucketObject {
	var tier StorageTierEnum
	switch o.StorageTier {
	case objectstorage.StorageTierStandard:
		tier = STierStandard
		break
	case objectstorage.StorageTierInfrequentAccess:
		tier = STierLowAccess
		break
	case objectstorage.StorageTierArchive:
		tier = STierTierArchive
		break
	}
	modfied := time.Now()

	if o.TimeModified != nil {
		modfied = o.TimeModified.Time
	}
	return BucketObject{
		Key:          *o.Name,
		LastModified: modfied,
		Size:         *o.Size,
		StorageClass: tier,
	}
}
