package main

import (
	"fmt"
	"testing"
	"time"
)

func TestConvertStringToTime(t *testing.T) {
	s := "2006-01-02T15:04:05.000Z"
	result := convertStringToTime(s)
	parsedStr, _ := time.Parse(s, s)

	if result != parsedStr {
		t.Errorf("Time conversion didn't match")
	}
}

func TestSortFilesByLastModified(t *testing.T) {
	obj1 := ObjectS3{
		Key:          "some_name1",
		LastModified: time.Now(),
	}
	obj2 := ObjectS3{
		Key:          "some_name2",
		LastModified: time.Now().Add(100),
	}
	objects := []ObjectS3{obj1, obj2}

	objPosition0 := objects[0]

	sortFilesByLastModified(objects)
	if objPosition0 == objects[0] {
		t.Errorf("Objects in slice were not sorted properly")
	}
}

func TestGetBucket(t *testing.T) {
	s := &S3{}
	s.Bucket = "some-bucket"
	s.AccessKey = "****"
	s.ClientSecret = "***"

	_, err := s.getBucket()
	if err != nil {
		t.Errorf("Bucket is not returned even if its name is valid")
	}
}

func TestUploadWithDumpThatHasError(t *testing.T) {
	s := &S3{}
	s.Bucket = "some-bucket"
	s.AccessKey = "****"
	s.ClientSecret = "***"

	dump := &ExportResult{}
	dump.Error = fmt.Errorf("test error")

	err := s.Upload(dump)
	if err == nil {
		t.Errorf("Dump with Error field passed to upload method must throw error")
	}
}

func TestUploadWithIncorrectBucketName(t *testing.T) {
	s := &S3{}
	s.Bucket = "some@@@bucket"
	s.AccessKey = "****"
	s.ClientSecret = "***"

	dump := &ExportResult{}
	err := s.Upload(dump)
	if err == nil {
		t.Errorf("Incorrect bucket name didn't throw and error")
	}
}
