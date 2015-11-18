package vault

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	//	"bytes"
)

func getVaultFs(algo string, key []byte) *Vault {
	root := "/tmp/goapp/test/vault"

	v := &Vault{
		Algo:    algo,
		BaseKey: key,
		Driver: &DriverFs{
			Root: root,
		},
	}

	os.RemoveAll(root)

	return v
}

func getVaultS3(algo string, key []byte) *Vault {
	root := os.Getenv("GONODE_TEST_AWS_VAULT_ROOT")

	if len(os.Getenv("TRAVIS_JOB_NUMBER")) > 0 {
		root += "/" + os.Getenv("TRAVIS_JOB_NUMBER")
	}

	if len(root) == 0 {
		root = "local"
	}

	bucket := os.Getenv("GONODE_TEST_AWS_VAULT_BUCKET")
	if len(bucket) == 0 {
		bucket = "gonode-test"
	}

	fmt.Printf("bucket: %s, root: %s\n", bucket, root)

	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{
			Filename: os.Getenv("HOME") + "/.aws/credentials",
			Profile:  "gonode-test",
		},
		&credentials.SharedCredentialsProvider{
			Filename: os.Getenv("GONODE_TEST_AWS_CREDENTIALS_FILE"),
			Profile:  os.Getenv("GONODE_TEST_AWS_PROFILE"),
		},
	})

	driver := &DriverS3{
		Root:        root,
		Region:      "eu-west-1",
		EndPoint:    "s3-eu-west-1.amazonaws.com",
		Bucket:      bucket,
		Credentials: creds,
	}

	v := &Vault{
		Algo:    algo,
		BaseKey: key,
		Driver:  driver,
	}

	driver.init()

	// delete objects
	l, _ := driver.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(root),
	})

	for _, o := range l.Contents {
		driver.client.DeleteObject(&s3.DeleteObjectInput{
			Key:    o.Key,
			Bucket: aws.String(bucket),
		})
	}

	return v
}

var algos = map[string][][]byte{
	"no_op":   {[]byte(""), key},
	"aes_ofb": {[]byte(""), key},
	"aes_ctr": {[]byte(""), key},
	"aes_cbc": {[]byte(""), key},
	"aes_gcm": {[]byte(""), key},
}

func runTest(driver string, t *testing.T, f func(algo string, key []byte) *Vault) {
	var m string
	for algo, keys := range algos {
		for _, key := range keys {
			v := f(algo, key)

			assert.False(t, v.Has("salut"), m+" - assert file does not exist")

			m = fmt.Sprintf("Type: %s/%s/xSmallMessage", driver, algo)
			t.Log(m)
			RunTestVault(t, v, xSmallMessage, m)

			m = fmt.Sprintf("Type: %s/%s/smallMessage", driver, algo)
			t.Log(m)
			RunTestVault(t, v, smallMessage, m)

			if _, travis := syscall.Getenv("TRAVIS"); travis == false {
				continue
			}

			m = fmt.Sprintf("Type: %s/%s/largeMessage", driver, algo)
			t.Log(m)
			RunTestVault(t, v, largeMessage, m)

			m = fmt.Sprintf("Type: %s/%s/xLargeMessage", driver, algo)
			t.Log(m)
			RunTestVault(t, v, xLargeMessage, m)
		}
	}
}

func Test_Vault_Drivers_FS(t *testing.T) {
	runTest("fs", t, getVaultFs)
}

func Test_Vault_Drivers_S3(t *testing.T) {
	if _, offline := syscall.Getenv("GONODE_TEST_OFFLINE"); offline == true {
		t.Skip("OFFLINE TEST ONLY")

		return
	}

	runTest("s3", t, getVaultS3)
}

//func Test_Generate_Regression_Files(t *testing.T) {
////	types := []string{"aes_ofb", "aes_ctr", "aes_cbc"}
//	types := []string{"aes_gcm"}
//
//	for _, v := range types {
//		v := &Vault{
//			Algo:    v,
//			BaseKey: []byte("de4d3ae8cf578c971b39ab5f21b2435483a3654f63b9f3777925c77e9492a141"),
//			Driver: &DriverFs{
//				Root: "../test/vault/" + v,
//			},
//		}
//
//		file := "The secret file"
//		data := bytes.NewBufferString("The secret message")
//		meta := NewVaultMetadata()
//		meta["foo"] = "bar"
//
//		if v.Has(file) {
//			v.Remove(file)
//		}
//
//		v.Put(file, meta, data)
//	}
//}

func getNoRegressionVaultFs(algo string) *Vault {
	return &Vault{
		Algo:    algo,
		BaseKey: key,
		Driver: &DriverFs{
			Root: "../test/vault/" + algo,
		},
	}
}

func Test_VaultFS_Secure_Aes_OFB_NoRegression(t *testing.T) {
	RunRegressionTest(t, getNoRegressionVaultFs("aes_ofb"))
}

func Test_VaultFS_Secure_Aes_CTR_NoRegression(t *testing.T) {
	RunRegressionTest(t, getNoRegressionVaultFs("aes_ctr"))
}

func Test_VaultFS_Secure_Aes_CBC_NoRegression(t *testing.T) {
	RunRegressionTest(t, getNoRegressionVaultFs("aes_cbc"))
}

func Test_VaultFS_Secure_Aes_GCM_NoRegression(t *testing.T) {
	RunRegressionTest(t, getNoRegressionVaultFs("aes_gcm"))
}
