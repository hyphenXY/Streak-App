package utils

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Uploader struct {
	BucketName string
	Client     *s3.S3
	Region     string
	Endpoint   string
}

func NewS3Uploader() (*S3Uploader, error) {
	region := os.Getenv("S3_REGION")
	endpoint := os.Getenv("S3_ENDPOINT") // works for Backblaze too
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET_KEY")
	bucket := os.Getenv("S3_BUCKET")

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	return &S3Uploader{
		BucketName: bucket,
		Client:     s3.New(sess),
		Region:     region,
		Endpoint:   endpoint,
	}, nil
}

// UploadFile uploads to directory like "profileimages/user/"
func (u *S3Uploader) UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, dir string) (string, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", err
	}

	fileKey := filepath.Join(dir, fileHeader.Filename)

	_, err := u.Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(fileKey),
		Body:   bytes.NewReader(buf.Bytes()),
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	fileURL := fmt.Sprintf("%s/%s/%s", u.Endpoint, u.BucketName, fileKey)
	return fileURL, nil
}
