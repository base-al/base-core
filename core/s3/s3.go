package s3

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	ORG_SETTINGS_FOLDER_NAME       = "settings"
	ORG_COMPANIES_FOLDER_NAME      = "companies"
	ORG_RESOURCES_FOLDER_NAME      = "resources"
	ORG_COMMUNITYPOSTS_FOLDER_NAME = "community-posts"
)

type S3Config struct {
	APIKey      string
	APISecret   string
	APIEndpoint string
	Region      string
	Env         string
	OwnerAccID  string
}

type BucketFile struct {
	FileName string
	Data     multipart.File
}

type BucketFileReply struct {
	EncBase64 string
}

type S3Client struct {
	c           *s3.S3
	env         string
	apiEndpoint string
	ownerAccId  string
}

func New(cfg *S3Config) *S3Client {
	// S3 New Session
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.APIKey, cfg.APISecret, ""),
		Endpoint:         &cfg.APIEndpoint,
		Region:           &cfg.Region,
		S3ForcePathStyle: aws.Bool(true),
	}))
	// service client
	s3c := s3.New(sess)
	return &S3Client{
		c:           s3c,
		env:         cfg.Env,
		apiEndpoint: cfg.APIEndpoint,
		ownerAccId:  cfg.OwnerAccID,
	}
}

func (s3c *S3Client) NewBucket(name *string) error {
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}
	// Creates new bucket
	_, err := s3c.c.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s3c *S3Client) NewOrgBucket(orgId int) error {
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}
	bucketOrgName := s3c.OrgBucketName(orgId)
	// Creates new bucket
	_, err := s3c.c.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: bucketOrgName,
	})
	if err != nil {
		return err
	}

	// this will let us to handle access using policies
	_, err = s3c.c.DeletePublicAccessBlockWithContext(ctx, &s3.DeletePublicAccessBlockInput{
		Bucket:              bucketOrgName,
		ExpectedBucketOwner: &s3c.ownerAccId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s3c *S3Client) UpdateOrgSettingsPublicPolicy(orgId int) error {
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}
	bucketOrgName := s3c.OrgBucketName(orgId)
	settingsTargetResource := fmt.Sprintf("arn:aws:s3:::%s/settings/*", *bucketOrgName)
	companiesTargetResource := fmt.Sprintf("arn:aws:s3:::%s/companies/*", *bucketOrgName)
	resourcesTargetResource := fmt.Sprintf("arn:aws:s3:::%s/resources/*", *bucketOrgName)
	communityPostsTargetResource := fmt.Sprintf("arn:aws:s3:::%s/community-posts/*", *bucketOrgName)

	readPolicy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "%s"
			},
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "%s"
			},
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "%s"
			},
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "%s"
			}
		]
	}`, settingsTargetResource, companiesTargetResource, resourcesTargetResource, communityPostsTargetResource)

	_, err := s3c.c.PutBucketPolicyWithContext(ctx, &s3.PutBucketPolicyInput{
		Bucket:              bucketOrgName,
		Policy:              &readPolicy,
		ExpectedBucketOwner: &s3c.ownerAccId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (s3c *S3Client) NewBucketFolder(bucketName string, folderName string) error {

	fmt.Println("NewBucketFolder", bucketName, folderName)
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}
	// Creates new folder in bucket
	_, err := s3c.c.PutObjectWithContext(ctx, &s3.PutObjectInput{

		Bucket: aws.String(bucketName),
		Key:    aws.String(folderName + "/"),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s3c *S3Client) Upload(bf *BucketFile, bucketName string) (resp *s3.PutObjectOutput, err error) {
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	// Uploads the object to S3. The Context will interrupt the request if the timeout expires
	resp, err = s3c.c.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(bf.FileName),
		Body:   bf.Data,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s3c *S3Client) ObjectURL(bucketName string, key string) string {
	prot := "https"
	return fmt.Sprintf("%s://%s/%s/%s", prot, s3c.apiEndpoint, bucketName, key)
}

func (s3c *S3Client) Download(fileName string, bucketName string) (resp *BucketFileReply, err error) {
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	s3Response, err := s3c.c.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}

	// Make sure to close the body when done with it for S3 GetObject APIs or will leak connections.
	defer s3Response.Body.Close()

	//fmt.Println("APIEndpoint", APIEndpoint)

	bytes, err := io.ReadAll(s3Response.Body)
	if err != nil {
		return nil, err
	}
	encBase64 := base64.StdEncoding.EncodeToString(bytes)
	return &BucketFileReply{
		EncBase64: encBase64,
	}, nil
}

// OrgBucketName combines bucket name in this format: 1.local, 1.dev, 1.prod
func (s3c *S3Client) OrgBucketName(orgId int) *string {
	bn := fmt.Sprintf("%d.org", orgId)
	if s3c.env == "" || s3c.env == "dev" {
		bn += ".local"
	}
	if s3c.env == "stage" {
		bn += ".stage"
	}
	if s3c.env == "test" {
		bn += ".test"
	}
	if s3c.env == "prod" {
		bn += ".prod"
	}
	return &bn
}

func (s3c *S3Client) UsersBucketName() string {
	bn := "users"
	if s3c.env == "" || s3c.env == "dev" {
		bn += ".local"
	}
	if s3c.env == "stage" {
		bn += ".stage"
	}
	if s3c.env == "test" {
		bn += ".test"
	}
	if s3c.env == "prod" {
		bn += ".prod"
	}
	return bn
}

// DeleteObject deletes the folder and its data inside
func (s3c *S3Client) DeleteBucketFolder(bucketName string, folderName string) error {
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	res, err := s3c.c.ListObjects(&s3.ListObjectsInput{
		Bucket: &bucketName,
	})
	if err != nil {
		fmt.Println("s3c.c.ListObjects", err)
	}

	dirPrefix := ""
	objIdentifiers := []*s3.ObjectIdentifier{}
	for _, obj := range res.Contents {
		dirPrefix = strings.Split(*obj.Key, "/")[0]
		if dirPrefix == folderName {
			objIdentifiers = append(objIdentifiers, &s3.ObjectIdentifier{
				Key: obj.Key,
			})
		}
	}
	if len(objIdentifiers) > 0 {
		_, err = s3c.c.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: &bucketName,
			Delete: &s3.Delete{
				Objects: objIdentifiers,
			},
		})
		if err != nil {
			fmt.Println("s3c.c.DeleteObjects", err)
			return err
		}
	}

	// Delete the "folder" object itself.
	_, err = s3c.c.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(folderName),
	})
	if err != nil {
		fmt.Println("Error deleting folder object:", err)
		return err
	}

	return nil
}

// DeleteBucket deletes the bucket and its data in S3
func (s3c *S3Client) DeleteBucket(name string) error {
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	listObjOut, err := s3c.c.ListObjects(&s3.ListObjectsInput{
		Bucket: &name,
	})
	if err != nil {
		fmt.Println("s3c.c.ListObjects", err)
	}
	objIdentifiers := []*s3.ObjectIdentifier{}
	for _, obj := range listObjOut.Contents {
		objIdentifiers = append(objIdentifiers, &s3.ObjectIdentifier{
			Key: obj.Key,
		})
	}

	if len(objIdentifiers) > 0 {
		_, err = s3c.c.DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: &name,
			Delete: &s3.Delete{
				Objects: objIdentifiers,
			},
		})
		if err != nil {
			fmt.Println("s3c.c.DeleteObjects", err)
			return err
		}
	}

	_, err = s3c.c.DeleteBucketWithContext(ctx, &s3.DeleteBucketInput{
		Bucket: &name,
	})
	if err != nil {
		fmt.Println("s3c.c.DeleteBucketWithContext.err", err)
		return err
	}
	return nil
}

func (s3c *S3Client) DeleteS3Object(bucketName string, key string) error {
	timeout := time.Duration(time.Minute * 10)
	ctx := context.Background()
	var cancelFn func()
	if timeout > 0 {
		ctx, cancelFn = context.WithTimeout(ctx, timeout)
	}
	if cancelFn != nil {
		defer cancelFn()
	}

	// Delete the "folder" object itself.
	_, err := s3c.c.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Println("Error deleting object:", err)
		return err
	}
	return nil
}
