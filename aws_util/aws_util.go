package aws_util

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"os"
	error2 "social-alarm-service/error"
)

type AWSUtil interface {
	UploadObject(ctx *gin.Context, fileName string, bucketName string, key string) *error2.ASError
	DeleteObject(ctx *gin.Context, bucketName string, key string) *error2.ASError
}

type awsUtil struct {
	s3Client *s3.Client
}

func NewAWSUtil(s3Client *s3.Client) AWSUtil {
	return awsUtil{s3Client: s3Client}
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
		ACL:    types.ObjectCannedACLPublicRead,
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
