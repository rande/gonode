// Copyright © 2014-2016 Thomas Rabaix <thomas.rabaix@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vault

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type s3Writer struct {
	client *s3.S3
	file   *os.File
	bucket string
	key    string
}

func (w *s3Writer) Write(b []byte) (int, error) {
	return w.file.Write(b)
}

func (w *s3Writer) Close() error {
	name := w.file.Name()

	defer func() {
		os.Remove(name)
	}()

	w.file.Seek(0, 0)

	_, err := w.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(w.bucket),
		Key:         aws.String(w.key),
		Body:        w.file,
		ContentType: aws.String("application/octet-stream"),
	})

	return err
}

type DriverS3 struct {
	Root        string
	Bucket      string
	Region      string
	EndPoint    string
	Credentials *credentials.Credentials
	client      *s3.S3
}

func (d *DriverS3) init() {
	if d.client != nil {
		return
	}

	config := &aws.Config{
		Region:           aws.String(d.Region),
		Endpoint:         aws.String(d.EndPoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      d.Credentials,
	}

	d.client = s3.New(session.New(), config)
}

func (d *DriverS3) getFilename(name string) string {
	return filepath.Join(d.Root, name)
}

func (d *DriverS3) Has(name string) bool {
	d.init()

	headResult, _ := d.client.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.getFilename(name)),
	})

	return headResult.ETag != nil
}

func (d *DriverS3) GetReader(name string) (io.ReadCloser, error) {
	d.init()

	result, err := d.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.getFilename(name)),
	})

	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (d *DriverS3) GetWriter(name string) (io.WriteCloser, error) {
	d.init()

	file, err := ioutil.TempFile(os.TempDir(), "s3_uploads_")

	if err != nil {
		return nil, err
	}

	return &s3Writer{
		bucket: d.Bucket,
		key:    d.getFilename(name),
		file:   file,
		client: d.client,
	}, nil
}

func (d *DriverS3) Remove(name string) error {
	d.init()

	_, err := d.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.getFilename(name)),
	})

	return err
}
