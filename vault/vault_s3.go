package vault

import (
//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/aws/awserr"
//	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
//	"github.com/aws/aws-sdk-go/service/s3"
)

type VaultS3 struct {
	Root        string
	Algo        string
	BaseKey     []byte
	Region      string
	EndPoint    string
	Credentials *credentials.Credentials
}

