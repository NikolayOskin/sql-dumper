package main

import (
	"bufio"
	"os"

	"gopkg.in/amz.v3/aws"
	"gopkg.in/amz.v3/s3"
)

// S3 is AWS S3 storer
type S3 struct {
	Region       string
	Bucket       string
	AccessKey    string
	ClientSecret string
}

// Upload dump sql file to storage
func (x *S3) Upload(dump *ExportResult) error {
	if dump.Error != nil {
		return dump.Error
	}

	auth := aws.Auth{
		AccessKey: x.AccessKey,
		SecretKey: x.ClientSecret,
	}
	s := s3.New(auth, aws.Regions[x.Region])

	bucket, err := s.Bucket(x.Bucket)
	if err != nil {
		return err
	}

	// Uploading
	file, err := os.Open(dump.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	buffy := bufio.NewReader(file)
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()

	err = bucket.PutReader(dump.Filename(), buffy, size, dump.MIME, s3.BucketOwnerFull)
	return err
}

func (x *S3) DeleteFile(filename string) error {
	auth := aws.Auth{
		AccessKey: x.AccessKey,
		SecretKey: x.ClientSecret,
	}
	s := s3.New(auth, aws.Regions[x.Region])

	bucket, err := s.Bucket(x.Bucket)
	if err != nil {
		return err
	}

	err = bucket.Del(filename)
	if err != nil {
		return err
	}

	return nil
}

func (x *S3) DeleteFilesExcept(filenames []string) error {
	auth := aws.Auth{
		AccessKey: x.AccessKey,
		SecretKey: x.ClientSecret,
	}
	s := s3.New(auth, aws.Regions[x.Region])

	bucket, err := s.Bucket(x.Bucket)
	if err != nil {
		return err
	}

	list, err := bucket.List("", "/", "", 1000)
	if err != nil {
		return err
	}

	for i := range list.Contents {
		if !stringInSlice(filenames, list.Contents[i].Key) {
			err = bucket.Del(list.Contents[i].Key)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func stringInSlice(slice []string, s string) bool {
	for _, b := range slice {
		if b == s {
			return true
		}
	}
	return false
}
