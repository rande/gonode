// Copyright © 2014-2023 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bytes"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getEnv(name, def string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		value = def
	}

	return value
}

func getChainCredentials() (*credentials.Credentials, error) {
	profile := getEnv("GONODE_TEST_AWS_PROFILE", "gonode-test")

	chainProvider := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{
			Filename: os.Getenv("HOME") + "/.aws/credentials",
			Profile:  profile,
		},
		&credentials.SharedCredentialsProvider{
			Filename: os.Getenv("GONODE_TEST_AWS_CREDENTIALS_FILE"),
			Profile:  profile,
		},
		&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     getEnv("GONODE_TEST_S3_ACCESS_KEY", ""),
			SecretAccessKey: getEnv("GONODE_TEST_S3_SECRET", ""),
		}},
	})

	if _, err := chainProvider.Get(); err != nil {
		return nil, err
	}

	return chainProvider, nil
}

func getDriver(chainProvider *credentials.Credentials) *DriverS3 {
	return &DriverS3{
		Bucket:      getEnv("GONODE_TEST_AWS_VAULT_S3_BUCKET", "gonode-qa"),
		Root:        getEnv("GITHUB_RUN_ID", getEnv("GONODE_TEST_AWS_VAULT_ROOT", "local")),
		Region:      getEnv("GONODE_TEST_S3_REGION", "eu-west-1"),
		EndPoint:    getEnv("GONODE_TEST_S3_ENDPOINT", "s3-eu-west-1.amazonaws.com"),
		Credentials: chainProvider,
	}
}

// this is just a test to validata how the aws sdk behave
func Test_Vault_Basic_S3_Usage(t *testing.T) {
	if getEnv("GONODE_TEST_OFFLINE", "yes") == "yes" {
		t.Skip("OFFLINE TEST ONLY")
		return
	}

	var err error
	var headResult *s3.HeadObjectOutput
	var getResult *s3.GetObjectOutput

	chainProvider, err := getChainCredentials()

	if err != nil {
		t.Skip("Unable to find credentials")
	}

	// init vault
	v := getDriver(chainProvider)

	// init credentials information
	config := &aws.Config{
		Region:           &v.Region,
		Endpoint:         &v.EndPoint,
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      v.Credentials,
	}

	s3client := s3.New(session.New(), config)

	key := fmt.Sprintf("%s/test/assd", v.Root)

	headResult, err = s3client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(v.Bucket),
		Key:    aws.String("no-file"),
	})

	assert.Error(t, err)
	assert.Nil(t, headResult.ETag)

	data := []byte("foobar et foo")

	putObject := &s3.PutObjectInput{
		Bucket:      aws.String(v.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("application/octet-stream"),
	}

	_, err = s3client.PutObject(putObject)

	headResult, err = s3client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(v.Bucket),
		Key:    aws.String(key),
	})

	assert.NoError(t, err)
	assert.NotNil(t, headResult.ETag)

	getObject := &s3.GetObjectInput{
		Bucket: aws.String(v.Bucket),
		Key:    aws.String(key),
	}

	getResult, err = s3client.GetObject(getObject)
	assert.NoError(t, err)

	data = []byte("xxxxxxxxxxxxx")

	getResult.Body.Read(data)
	getResult.Body.Close()

	assert.Equal(t, []byte("foobar et foo"), data)
}
