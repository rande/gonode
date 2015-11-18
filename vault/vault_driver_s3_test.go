// Copyright © 2014-2015 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"syscall"
)

// this is just a test to validata how the aws sdk behave
func Test_Vault_Basic_S3_Usage(t *testing.T) {

	if _, offline := syscall.Getenv("GONODE_TEST_OFFLINE"); offline == true {
		t.Skip("OFFLINE TEST ONLY")
		return
	}

	var err error
	var headResult *s3.HeadObjectOutput
	var getResult *s3.GetObjectOutput

	root := os.Getenv("GONODE_TEST_AWS_VAULT_ROOT")
	if len(root) == 0 {
		root = "local"
	}

	profile := os.Getenv("GONODE_TEST_AWS_PROFILE")
	if len(profile) == 0 {
		profile = "gonode-test"
	}

	// init vault
	v := &DriverS3{
		Root:     root,
		Region:   "eu-west-1",
		EndPoint: "s3-eu-west-1.amazonaws.com",
		Credentials: credentials.NewChainCredentials([]credentials.Provider{
			&credentials.EnvProvider{},
			&credentials.SharedCredentialsProvider{
				Filename: os.Getenv("HOME") + "/.aws/credentials",
				Profile:  profile,
			},
			&credentials.SharedCredentialsProvider{
				Filename: os.Getenv("GONODE_TEST_AWS_CREDENTIALS_FILE"),
				Profile:  profile,
			},
		}),
	}

	// init credentials information
	config := &aws.Config{
		Region:           &v.Region,
		Endpoint:         &v.EndPoint,
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      v.Credentials,
	}

	s3client := s3.New(session.New(), config)

	bucketName := os.Getenv("GONODE_TEST_AWS_VAULT_S3_BUCKET")
	if len(bucketName) == 0 {
		bucketName = "gonode-test"
	}

	key := fmt.Sprintf("%s/test/assd", v.Root)

	headResult, err = s3client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String("no-file"),
	})

	assert.Error(t, err)
	assert.Nil(t, headResult.ETag)

	data := []byte("foobar et foo")

	putObject := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/octet-stream"),
	}

	_, err = s3client.PutObject(putObject)

	headResult, err = s3client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})

	assert.NoError(t, err)
	assert.NotNil(t, headResult.ETag)

	getObject := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	getResult, err = s3client.GetObject(getObject)
	assert.NoError(t, err)

	data = []byte("xxxxxxxxxxxxx")

	getResult.Body.Read(data)
	getResult.Body.Close()

	assert.Equal(t, []byte("foobar et foo"), data)
}
