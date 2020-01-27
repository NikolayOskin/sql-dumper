package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"

	"gopkg.in/amz.v3/aws"
	"gopkg.in/amz.v3/s3"
)

// S3 is AWS S3 storer
type S3 struct {
	Region       string
	BucketName   string
	AccessKey    string
	ClientSecret string
	Bucket       Bucket
}

type ObjectS3 struct {
	Key          string
	LastModified time.Time
}

type Bucket interface {
	List(prefix string, delim string, marker string, max int) (*s3.ListResp, error)
	Del(name string) error
	PutReader(path string, r io.Reader, length int64, contType string, perm s3.ACL) error
}

// Upload dump sql file to storage
func (x *S3) Upload(dump *ExportResult) error {
	if dump.Error != nil {
		return dump.Error
	}

	bucket, err := x.getBucket()
	if err != nil {
		return err
	}

	// Uploading
	file, err := os.Open(dump.Path)
	if err != nil {
		return fmt.Errorf("failed open file %s: %s", dump.Path, err)
	}

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to return stat for %s: %s", dump.Path, err)
	}

	size := stat.Size()

	err = bucket.PutReader(dump.Filename(), file, size, dump.MIME, s3.BucketOwnerFull)
	if err != nil {
		file.Close()
		return fmt.Errorf("failed to upload file %s to S3: %s", dump.Filename(), err)
	}

	file.Close()

	return nil
}

func (x *S3) DeleteOldFiles(dumpsToKeep int) error {
	bucket, err := x.getBucket()
	if err != nil {
		return err
	}

	list, err := bucket.List("", "/", "", 1000)
	if err != nil {
		return fmt.Errorf("trying to list object from S3 bucket: %v", err)
	}

	objects := x.getObjects(list)

	sortByLastModified(objects)
	latestObjects := keepLatest(objects, dumpsToKeep)

	for i := range list.Contents {
		if !isLatest(list.Contents[i].Key, latestObjects) {
			err := bucket.Del(list.Contents[i].Key)
			if err != nil {
				log.Printf("failed to delete object %s from S3 bucket: %s", list.Contents[i].Key, err)
			}
		}
	}

	return nil
}

func (x *S3) getObjects(list *s3.ListResp) []ObjectS3 {
	var objects []ObjectS3

	for i := range list.Contents {
		object := ObjectS3{
			Key:          list.Contents[i].Key,
			LastModified: convertStringToTime(list.Contents[i].LastModified),
		}
		objects = append(objects, object)
	}
	return objects
}

func (x *S3) getBucket() (Bucket, error) {
	if x.Bucket != nil {
		return x.Bucket, nil
	}

	auth := aws.Auth{
		AccessKey: x.AccessKey,
		SecretKey: x.ClientSecret,
	}
	s := s3.New(auth, aws.Regions[x.Region])

	bucket, err := s.Bucket(x.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to return S3 bucket: %v", err)
	}

	x.Bucket = bucket

	return x.Bucket, nil
}

func isLatest(key string, latestObjects []ObjectS3) bool {
	for _, object := range latestObjects {
		if object.Key == key {
			return true
		}
	}
	return false
}

func keepLatest(objects []ObjectS3, dumpsToKeep int) []ObjectS3 {
	if len(objects) > dumpsToKeep {
		latestObjects := objects[:dumpsToKeep]
		return latestObjects
	}
	return objects
}

func convertStringToTime(str string) time.Time {
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}

	return t
}

func sortByLastModified(objects []ObjectS3) {
	sort.SliceStable(objects, func(i, j int) bool {
		return objects[i].LastModified.After(objects[j].LastModified)
	})
}
