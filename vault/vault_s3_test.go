package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
//	"github.com/aws/aws-sdk-go/aws/awserr"
//	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"fmt"
//	"time"
	"bytes"
	"os"
)

func Test_Vault_Basic_S3_Usage(t *testing.T) {
	var err error
	var headResult *s3.HeadObjectOutput

	// init vault
	v := &VaultS3{
		Region: "eu-west-1",
		EndPoint: "s3-eu-west-1.amazonaws.com",
		Credentials: credentials.NewChainCredentials([]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Filename: os.Getenv("HOME") + "/.aws/credentials",
				Profile: "gonode-test",
			},
			&credentials.SharedCredentialsProvider{
				Filename: os.Getenv("HOME") + "/.aws/credentials",
				Profile: os.Getenv("GONODE_TEST_AWS_PROFILE"),
			},
		}),
	}

	// init credentials information
	config := &aws.Config{
		Region:           &v.Region,
		Endpoint:         &v.EndPoint, // <-- forking important !
		S3ForcePathStyle: aws.Bool(true), // <-- without these lines. All will fail! fork you aws!
		Credentials:      v.Credentials,
	}

	s3client := s3.New(session.New(), config)

	bucketName := os.Getenv("GONODE_TEST_AWS_VAULT_S3_BUCKET")
	if len(bucketName) == 0 {
		bucketName = "gonode-test"
	}

	key := fmt.Sprintf("/test/assd")

	headResult, err = s3client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String("/no-file"),
	})

	assert.Error(t, err)
	assert.Nil(t, headResult.ETag)

	//	now := time.Now()
	//	key := fmt.Sprintf("/test/%s", now.String())
	data := []byte("foobar et foo")

	putObject := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName), // required
		Key:           aws.String(key), // required
		Body:          bytes.NewReader(data),
		ContentType:   aws.String("application/octet-stream"),
	}

	_, err = s3client.PutObject(putObject)

	headResult, err = s3client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(key),
	})

	assert.NoError(t, err)
	assert.NotNil(t, headResult.ETag)

	assert.NoError(t, err)

	//	v := &VaultS3{
	//		Root: "/",
	//		Algo: "aes_ctr",
	//		BaseKey: key,
	//		Region: "eu-west-1",
	//		Credentials: creds,
}
