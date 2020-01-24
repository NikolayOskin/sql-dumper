package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"

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

type ObjectS3 struct {
	Key          string
	LastModified time.Time
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
		return fmt.Errorf("trying to return S3 bucket when uploading: %v", err)
	}

	// Uploading
	file, err := os.Open(dump.Path)
	if err != nil {
		return fmt.Errorf("trying to open dump file when uploading to S3: %v", err)
	}

	defer file.Close()

	buffy := bufio.NewReader(file)
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("trying to get dump file stat when uploading to S3: %v", err)
	}

	size := stat.Size()

	err = bucket.PutReader(dump.Filename(), buffy, size, dump.MIME, s3.BucketOwnerFull)
	if err != nil {
		return fmt.Errorf("trying to insert object to S3: %v", err)
	}

	return nil
}

func (x *S3) DeleteFilesExceptLatest(dumpsToKeep uint) error {
	auth := aws.Auth{
		AccessKey: x.AccessKey,
		SecretKey: x.ClientSecret,
	}
	s := s3.New(auth, aws.Regions[x.Region])

	bucket, err := s.Bucket(x.Bucket)
	if err != nil {
		return fmt.Errorf("trying return S3 bucket: %v", err)
	}

	list, err := bucket.List("", "/", "", 1000)
	if err != nil {
		return fmt.Errorf("trying to list object from S3 bucket: %v", err)
	}

	var objects []ObjectS3

	for i := range list.Contents {
		object := ObjectS3{
			Key:          list.Contents[i].Key,
			LastModified: convertStringToTime(list.Contents[i].LastModified),
		}
		objects = append(objects, object)

	}

	sortByDate(objects)

	objects = objects[:dumpsToKeep]

	for _, object := range objects {
		err = bucket.Del(object.Key)
		if err != nil {
			return fmt.Errorf("trying to delete object from S3 bucket: %v", err)
		}
	}

	return nil
}

func convertStringToTime(str string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}

	return t
}

func sortByDate(objects []ObjectS3) {
	sort.SliceStable(objects, func(i, j int) bool {
		return objects[i].LastModified.After(objects[j].LastModified)
	})
}
