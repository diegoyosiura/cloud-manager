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
	lastModified := time.Now()
	key := ""
	size := int64(0)
	if o.LastModified != nil {
		lastModified = *o.LastModified
	}

	if o.Key != nil {
		key = *o.Key
	}
	if o.Size != nil {
		size = *o.Size
	}

	return BucketObject{
		Key:          key,
		LastModified: lastModified,
		Size:         size,
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
	lastModified := time.Now()
	key := ""
	size := int64(0)
	if o.TimeModified != nil {
		lastModified = o.TimeModified.Time
	}

	if o.Name != nil {
		key = *o.Name
	}
	if o.Size != nil {
		size = *o.Size
	}

	return BucketObject{
		Key:          key,
		LastModified: lastModified,
		Size:         size,
		StorageClass: tier,
	}
}
