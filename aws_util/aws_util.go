package aws_util

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	error2 "social-alarm-service/error"
)

type AWSUtil interface {
	UploadObject(ctx *gin.Context, file *os.File, bucketName string, key string) *error2.ASError
}

type awsUtil struct {
	s3Client *s3.Client
}

func NewAWSUtil(s3Client *s3.Client) AWSUtil {
	return awsUtil{s3Client: s3Client}
}

func (awsUtil awsUtil) UploadObject(ctx *gin.Context, file *os.File, bucketName string, key string) *error2.ASError {
	_, err := awsUtil.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Printf("Couldn't upload file %v to %v:%v. Here's why: %v\n",
			file.Name(), bucketName, key, err)
		return error2.InternalServerError("aws upload object failed")
	}
	return nil
}
