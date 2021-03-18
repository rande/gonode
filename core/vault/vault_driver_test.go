// Copyright © 2014-2021 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	//		"bytes"
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

	creds, err := getChainCredentials()

	if err != nil {
		return nil
	}

	driver := getDriver(creds)

	v := &Vault{
		Algo:    algo,
		BaseKey: key,
		Driver:  driver,
	}

	driver.init()

	// delete objects
	l, _ := driver.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(driver.Bucket),
		Prefix: aws.String(driver.Root),
	})

	for _, o := range l.Contents {
		fmt.Printf("Delete: %s / %s\n", driver.Bucket, *o.Key)
		driver.client.DeleteObject(&s3.DeleteObjectInput{
			Key:    o.Key,
			Bucket: aws.String(driver.Bucket),
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

			if v == nil {
				t.Skip("Unable to get vault (missing credentials?)")

				return
			}

			assert.False(t, v.Has("salut"), m+" - assert file does not exist")

			m = fmt.Sprintf("Type: %s/%s/xSmallMessage", driver, algo)
			t.Log(m)
			RunTestVault(t, v, xSmallMessage, m)

			m = fmt.Sprintf("Type: %s/%s/smallMessage", driver, algo)
			t.Log(m)
			RunTestVault(t, v, smallMessage, m)

			// m = fmt.Sprintf("Type: %s/%s/largeMessage", driver, algo)
			// t.Log(m)
			// RunTestVault(t, v, largeMessage, m)

			// m = fmt.Sprintf("Type: %s/%s/xLargeMessage", driver, algo)
			// t.Log(m)
			// RunTestVault(t, v, xLargeMessage, m)
		}
	}
}

func Test_Vault_Drivers_FS(t *testing.T) {
	//runTest("fs", t, getVaultFs)
}

func Test_Vault_Drivers_S3(t *testing.T) {
	if getEnv("GONODE_TEST_OFFLINE", "yes") == "yes" {
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
			Root: "../../test/vault/" + algo,
		},
	}
}

// func Test_VaultFS_Secure_Aes_OFB_NoRegression(t *testing.T) {
// 	RunRegressionTest(t, getNoRegressionVaultFs("aes_ofb"))
// }

// func Test_VaultFS_Secure_Aes_CTR_NoRegression(t *testing.T) {
// 	RunRegressionTest(t, getNoRegressionVaultFs("aes_ctr"))
// }

// func Test_VaultFS_Secure_Aes_CBC_NoRegression(t *testing.T) {
// 	RunRegressionTest(t, getNoRegressionVaultFs("aes_cbc"))
// }

// func Test_VaultFS_Secure_Aes_GCM_NoRegression(t *testing.T) {
// 	RunRegressionTest(t, getNoRegressionVaultFs("aes_gcm"))
// }
