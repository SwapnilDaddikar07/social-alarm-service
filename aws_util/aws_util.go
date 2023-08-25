package aws_util

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	error2 "social-alarm-service/error"
	"time"
)

type AWSUtil interface {
	UploadObject(ctx *gin.Context, fileName string, bucketName string, key string) *error2.ASError
	DeleteObject(ctx *gin.Context, bucketName string, key string) *error2.ASError
	GeneratePresignedURL(ctx *gin.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, *error2.ASError)
}

type awsUtil struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
}

func NewAWSUtil(s3Client *s3.Client, presignClient *s3.PresignClient) AWSUtil {
	return awsUtil{s3Client: s3Client, presignClient: presignClient}
}

func (awsUtil awsUtil) UploadObject(ctx *gin.Context, fileName string, bucketName string, key string) *error2.ASError {
	fmt.Printf("opening file %s\n", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("could not open file %s\n", fileName)
		return error2.InternalServerError("could not open tmp file")
	}
	defer file.Close()

	fmt.Println("file opened successfully. uploading object to s3")
	_, err = awsUtil.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
		ACL:    types.ObjectCannedACLPrivate,
	})
	if err != nil {
		fmt.Printf("Couldn't upload file to %v:%v. Here's why: %v\n", bucketName, key, err)
		return error2.InternalServerError("aws upload object failed")
	}

	fmt.Printf("file uploaded successfully to s3")
	return nil
}

func (awsUtil awsUtil) DeleteObject(ctx *gin.Context, bucketName string, key string) *error2.ASError {
	fmt.Printf("deleting %s key from bucket %s\n", key, bucketName)

	_, err := awsUtil.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	})
	if err != nil {
		fmt.Printf("could not delete key %s from bucket %s \n", key, bucketName)
		return error2.InternalServerError("could not delete object from s3")
	}

	fmt.Printf("successfully deleted %s key from bucket %s", key, bucketName)
	return nil
}

func (awsUtil awsUtil) GeneratePresignedURL(ctx *gin.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, *error2.ASError) {
	request, err := awsUtil.presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
		return nil, error2.InternalServerError("could not generate presigned URL")
	}

	fmt.Println(fmt.Sprintf("presigned url successfully generated for key %s", objectKey))
	return request, nil
}
