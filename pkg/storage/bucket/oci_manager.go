package bucket

import (
	"context"
	"fmt"
	"github.com/diegoyosiura/cloud-manager/pkg/authentication"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"github.com/oracle/oci-go-sdk/v65/objectstorage/transfer"
	"os"
	"time"
)

type OCIManager struct {
	Auth   *authentication.OCIAuth // OCI authentication details.
	Client *objectstorage.ObjectStorageClient
}

func (o *OCIManager) setup() (bool, error) {
	if o.Client == nil {
		c, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(o.Auth.GetConfigurationProvider())
		if err != nil {
			return false, err
		}

		o.Client = &c
	}

	return true, nil
}

func (o *OCIManager) List(name string) (r []BucketObject, err error) {
	successs, err := o.setup()
	if !successs {
		panic(err)
	}
	ctx := context.Background()
	rq := objectstorage.ListObjectsRequest{}

	rq.NamespaceName = &o.Auth.Namespace
	rq.BucketName = &name

	resp, err := o.Client.ListObjects(ctx, rq)
	if err != nil {
		return nil, err
	}

	for _, obj := range resp.ListObjects.Objects {
		r = append(r, NewBucketObjectFromOCI(obj))
	}

	return r, nil
}

func (o *OCIManager) Create(name string, waitCreate bool) error {
	successs, err := o.setup()
	if !successs {
		panic(err)
	}

	ctx := context.Background()
	rq := objectstorage.CreateBucketRequest{
		NamespaceName: &o.Auth.Namespace,
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			Name:          &name,
			CompartmentId: &o.Auth.CompartmentID,
		},
	}
	_, err = o.Client.CreateBucket(ctx, rq)

	if err != nil {
		return err
	}

	if waitCreate {
		for _, err = o.List(name); err != nil; {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (o *OCIManager) Delete(name string) error {
	successs, err := o.setup()
	if !successs {
		panic(err)
	}

	ctx := context.Background()
	rq := objectstorage.DeleteBucketRequest{
		NamespaceName: &o.Auth.Namespace,
		BucketName:    &name,
	}
	_, err = o.Client.DeleteBucket(ctx, rq)

	if err != nil {
		return err
	}

	return nil
}

func (o *OCIManager) Upload(bucket string, objectName string, f *os.File, partSize int64, threads int) error {
	successs, err := o.setup()
	if !successs {
		panic(err)
	}

	if partSize < 131072 { // 128 * 1024
		partSize = 10 * 1024 * 1024
	}
	if threads <= 0 { // 128 * 1024
		threads = 4
	}

	trueBool := true
	rq := transfer.UploadStreamRequest{
		UploadRequest: transfer.UploadRequest{
			NamespaceName:         &o.Auth.Namespace,
			BucketName:            &bucket,
			ObjectName:            &objectName,
			PartSize:              &partSize,
			AllowMultipartUploads: &trueBool,
			AllowParrallelUploads: &trueBool,
			NumberOfGoroutines:    &threads,
			ObjectStorageClient:   o.Client,
			StorageTier:           "STANDARD",
		},
		StreamReader: f,
	}
	uploader := transfer.NewUploadManager()

	ctx := context.Background()
	_, err = uploader.UploadStream(ctx, rq)

	if err != nil {
		return err
	}

	return nil
}

func (o *OCIManager) Update(bucket string, objectName string, f *os.File, partSize int64, threads int) error {
	return o.Upload(bucket, objectName, f, partSize, threads)
}

func (o *OCIManager) DownloadLink(bucketName string, objectName string, expires int64) (string, error) {
	successs, err := o.setup()
	if !successs {
		panic(err)
	}
	ctx := context.Background()

	expiration := common.SDKTime{Time: time.Now().Add(time.Duration(expires) * time.Minute)}
	rq := objectstorage.CreatePreauthenticatedRequestRequest{
		NamespaceName: &o.Auth.Namespace,
		BucketName:    &bucketName,
		CreatePreauthenticatedRequestDetails: objectstorage.CreatePreauthenticatedRequestDetails{
			Name:        common.String("temp-link-" + time.Now().Format("20060102150405")),
			AccessType:  objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectread,
			TimeExpires: &expiration,
			ObjectName:  &objectName,
		},
	}

	resp, err := o.Client.CreatePreauthenticatedRequest(ctx, rq)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://objectstorage.%s.oraclecloud.com%s", o.Auth.Region, *resp.PreauthenticatedRequest.AccessUri), nil
}
func (o *OCIManager) DeleteObject(bucketName string, objectName string) error {
	successs, err := o.setup()
	if !successs {
		panic(err)
	}
	ctx := context.Background()

	rq := objectstorage.DeleteObjectRequest{
		NamespaceName: &o.Auth.Namespace,
		BucketName:    &bucketName,
		ObjectName:    &objectName,
	}

	_, err = o.Client.DeleteObject(ctx, rq)
	if err != nil {
		return err
	}
	return nil
}
