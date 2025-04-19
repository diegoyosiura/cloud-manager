package bucket

import (
	"bytes"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/diegoyosiura/cloud-manager/pkg/authentication"
	"io"
	"os"
	"sort"
	"time"
)

type AWSManager struct {
	Auth   *authentication.AWSAuth // AWS authentication details.
	Client *s3.S3
}

func (a *AWSManager) setup() (bool, error) {
	if a.Client == nil {
		a.Client = s3.New(a.Auth.Session, &aws.Config{Region: &a.Auth.Region})
		if a.Client == nil {
			return false, errors.New("failed to create AWS client")
		}
	}

	return true, nil
}
func (a *AWSManager) List(name string) (r []BucketObject, err error) {
	successs, err := a.setup()
	if !successs {
		panic(err)
	}

	bi := &s3.ListObjectsV2Input{}
	bi.Bucket = &name

	buckets, err := a.Client.ListObjectsV2(bi)

	for _, b := range buckets.Contents {
		r = append(r, NewBucketObjectFromAWS(b))
	}

	return r, nil
}

func (a *AWSManager) Create(name string, waitCreate bool) error {
	successs, err := a.setup()
	if !successs {
		panic(err)
	}

	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err = a.Client.CreateBucket(input)

	if err != nil {
		return err
	}

	if waitCreate {
		err = a.Client.WaitUntilBucketExists(&s3.HeadBucketInput{
			Bucket: aws.String(name),
		})
	}
	return nil
}

func (a *AWSManager) Delete(name string) error {
	successs, err := a.setup()
	if !successs {
		panic(err)
	}

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}

	_, err = a.Client.DeleteBucket(input)

	if err != nil {
		return err
	}
	return nil
}

func (a *AWSManager) Upload(bucket string, objectName string, f *os.File, partSize int64, threads int) error {
	successs, err := a.setup()
	if !successs {
		panic(err)
	}

	if partSize < 131072 { // 128 * 1024
		partSize = 10 * 1024 * 1024
	}
	if threads <= 0 { // 128 * 1024
		threads = 4
	}

	rq := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectName),
	}

	initOut, err := a.Client.CreateMultipartUpload(rq)
	if err != nil {
		return err
	}

	uploadID := initOut.UploadId
	partNum := int64(1)
	buf := make([]byte, partSize)

	var completed []*s3.CompletedPart
	for {
		n, readErr := f.Read(buf)
		if n > 0 {
			out, err := a.upload(bucket, objectName, partNum, uploadID, buf, n)
			if err != nil {
				_, _ = a.Client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
					Bucket: aws.String(bucket), Key: aws.String(objectName), UploadId: uploadID,
				})
				return err
			}

			completed = append(completed, &s3.CompletedPart{
				ETag: out.ETag, PartNumber: aws.Int64(partNum),
			})
			partNum++
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	sort.Slice(completed, func(i, j int) bool {
		if *completed[i].PartNumber == *completed[j].PartNumber {
			return true
		}
		return *completed[i].PartNumber < *completed[j].PartNumber
	})
	_, err = a.Client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(objectName),
		UploadId: uploadID,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completed,
		},
	})
	return err
}

func (a *AWSManager) upload(bucket, objectName string, partNum int64, uploadID *string, buf []byte, n int) (*s3.UploadPartOutput, error) {
	out, err := a.Client.UploadPart(&s3.UploadPartInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(objectName),
		PartNumber: &partNum,
		UploadId:   uploadID,
		Body:       bytes.NewReader(buf[:n]),
	})
	if err != nil {
		_, _ = a.Client.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
			Bucket: aws.String(bucket), Key: aws.String(objectName), UploadId: uploadID,
		})
		return nil, err
	}

	return out, nil
}

func (a *AWSManager) Update(bucket string, objectName string, f *os.File, partSize int64, threads int) error {
	successs, err := a.setup()
	if !successs {
		panic(err)
	}

	return a.Upload(bucket, objectName, f, partSize, threads)
}
func (a *AWSManager) DownloadLink(bucketName string, objectName string, expires int64) (string, error) {
	successs, err := a.setup()
	if !successs {
		panic(err)
	}

	req, _ := a.Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	urlStr, err := req.Presign(time.Duration(expires) * time.Minute)

	if err != nil {
		return "", err
	}

	return urlStr, nil
}
func (a *AWSManager) DeleteObject(bucketName string, objectName string) error {
	successs, err := a.setup()
	if !successs {
		panic(err)
	}

	req := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	_, err = a.Client.DeleteObject(req)

	if err != nil {
		return err
	}
	return nil
}
