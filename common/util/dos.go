package util

import (
	"fmt"
	"git.multiverse.io/eventkit/kit/common/errors"
	"git.multiverse.io/eventkit/kit/constant"
	"git.multiverse.io/eventkit/kit/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"os"
	"time"
)

// CephMgmt is the configs of CEPH management
type CephMgmt struct {
	Host      string `ceph host`
	Region    string `ceph region`
	Bucket    string `bucket_id`
	AccessKey string `aws s3 aceessKey`
	SecretKey string `aws s3 secretKey`
	PathStyle bool   `ceph url style, true means you can use host directly,  false means bucket_id.doname will be used`
}

// Ceph is global variable to hold CEPH management connection
var Ceph CephMgmt

// MyProvider is an implement of Provider
type MyProvider struct{}

// Retrieve returns the credentials value from provider
func (m *MyProvider) Retrieve() (credentials.Value, error) {

	return credentials.Value{
		AccessKeyID:     Ceph.AccessKey,
		SecretAccessKey: Ceph.SecretKey,
	}, nil
}

// IsExpired never expired.
func (m *MyProvider) IsExpired() bool { return false }

// DosConfig for simplify the input of CEPH management connection parameters
type DosConfig struct {
	Host      string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}

// Init initializes the CEPH manager using DosConfig
func (c *CephMgmt) Init(config *DosConfig) error {

	c.Host = config.Host
	if c.Host == "" {
		return errors.New(constant.SystemInternalError, "no dos::host in config")
	}

	c.Region = config.Region

	c.Bucket = config.Bucket
	if c.Bucket == "" {
		return errors.New(constant.SystemInternalError, "no dos::bucket in config")
	}

	c.AccessKey = config.AccessKey
	if c.AccessKey == "" {
		return errors.New(constant.SystemInternalError, "no dos::accesskey in config")
	}

	c.SecretKey = config.SecretKey
	if c.SecretKey == "" {
		return errors.New(constant.SystemInternalError, "no dos::secretkey in config")
	}

	c.PathStyle = true

	return nil
}

// InitByParams initializes the CEPH manager using parameters
func (c *CephMgmt) InitByParams(host, region, bucket, accessKey, secretKey string) error {

	c.Host = host
	if c.Host == "" {
		return errors.New(constant.SystemInternalError, "no ceph::host in config")
	}

	c.Region = region
	if len(c.Region) == 0 {
		c.Region = "default"
	}

	c.Bucket = bucket
	if c.Bucket == "" {
		return errors.New(constant.SystemInternalError, "no ceph::bucket in config")
	}

	c.AccessKey = accessKey
	if c.AccessKey == "" {
		return errors.New(constant.SystemInternalError, "no ceph::accesskey in config")
	}

	c.SecretKey = secretKey
	if c.SecretKey == "" {
		return errors.New(constant.SystemInternalError, "no ceph::secretkey in config")
	}

	c.PathStyle = true

	return nil
}

func (c *CephMgmt) connect() (*s3.S3, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(c.Region),
			//EndpointResolver: endpoints.ResolverFunc(s3CustResolverFn),
			Endpoint:         &c.Host,
			S3ForcePathStyle: &c.PathStyle,
			Credentials:      credentials.NewCredentials(&MyProvider{}),
		},
	}))
	// Create the S3 service client with the shared session. This will
	// automatically use the S3 custom endpoint configured in the custom
	// endpoint resolver wrapping the default endpoint resolver.
	return s3.New(sess), nil
}

// UploadLargeFile upload large file to CEPH server
func (c *CephMgmt) UploadLargeFile(srcFile, cephFileKey string) error {
	s3Svc, _ := c.connect()
	uploader := s3manager.NewUploaderWithClient(s3Svc)

	f, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer f.Close()

	start := time.Now().Unix()
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(cephFileKey),
		Body:   f,
	})
	if err != nil {
		return err
	}
	log.Infosf("upload %s spend %ds", srcFile, time.Now().Unix()-start)
	return nil
}

// UploadFile upload file to CEPH server
func (c *CephMgmt) UploadFile(file io.Reader, cephFileKey string) (versionID string, err error) {
	s3Svc, _ := c.connect()
	uploader := s3manager.NewUploaderWithClient(s3Svc)
	output, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(cephFileKey),
		Body:   file,
	})
	if err != nil {
		return
	}
	if output.VersionID == nil {
		return "", nil
	}
	return *output.VersionID, nil
}

// DownloadFile download file from CEPH server
func (c *CephMgmt) DownloadFile(fileName, versionID string) (io.ReadCloser, error) {
	s3Svc, _ := c.connect()
	input := &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(fileName),
	}
	if len(versionID) > 0 {
		input.VersionId = aws.String(versionID)
	}
	output, err := s3Svc.GetObject(input)
	if err != nil {
		return nil, err
	}
	return output.Body, err
}

// DownloadLargeFile download large file from CEPH server
func (c *CephMgmt) DownloadLargeFile(cephFileKey, dstName string) error {
	s3Svc, _ := c.connect()
	downloader := s3manager.NewDownloaderWithClient(s3Svc)

	f, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer f.Close()

	start := time.Now().Unix()

	n, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(cephFileKey),
	})

	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "failed to download file, %v", err)
	}
	log.Infosf("download file %s , total %d bytes , spend %ds", dstName, n, time.Now().Unix()-start)
	return nil
}

// ListBuckets lists all the buckets in the CEPH server
func (c *CephMgmt) ListBuckets() {
	s3Svc, _ := c.connect()
	input := &s3.ListBucketsInput{}
	output, err := s3Svc.ListBuckets(input)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Printf("list %v", output)
}

// ListObjects lists all the objects in the CEPH server
func (c *CephMgmt) ListObjects(bucket string) {
	s3Svc, _ := c.connect()
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}
	output, err := s3Svc.ListObjects(input)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	log.Infosf("list %v", output)
}

// CreateBucket creates the bucket in the CEPH server
func (c *CephMgmt) CreateBucket(bucket string) error {
	s3Svc, _ := c.connect()

	crParams := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(""),
		},
	}

	_, err := s3Svc.CreateBucket(crParams)

	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Unable to create bucket %q, %v", bucket, err)
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)

	err = s3Svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return errors.Errorf(constant.SystemInternalError, "Error occurred while waiting for bucket[%s] to be created , %v", bucket, err)
	}

	fmt.Printf("Bucket %q successfully created\n", bucket)

	puParams := &s3.PutBucketAclInput{
		Bucket: aws.String(bucket),
	}
	puParams.SetACL("public-read") //set bucket ACL

	_, err = s3Svc.PutBucketAcl(puParams)
	if err != nil {
		return err
	}
	log.Infosf("Set", bucket, "ACL to public-read")
	return nil
}

// DeleteFile deletes the file in the CEPH server with file ID
func (c *CephMgmt) DeleteFile(fileID string) error {
	s3Svc, _ := c.connect()
	_, err := s3Svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(fileID),
	})
	return err
}
