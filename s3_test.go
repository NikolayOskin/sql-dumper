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

	firstObjBeforeSorting := objects[0]
	sortByLastModified(objects)
	firstObjAfterSorting := objects[0]

	if firstObjBeforeSorting == firstObjAfterSorting {
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

func TestKeepLatest(t *testing.T) {
	objects := generateObjects()

	// DumpsToKeep is less than slice length
	latestObjects := keepLatest(objects, 3)
	if len(latestObjects) != 3 {
		t.Errorf("Objects were not filtered")
	}

	// DumpsToKeep is more than slice length
	latestObjects = keepLatest(objects, 20)
	if len(latestObjects) != 4 {
		t.Errorf("Slice size must be 4")
	}
}

func generateObjects() []ObjectS3 {
	objects := []ObjectS3{
		{
			Key:          "some",
			LastModified: time.Now(),
		},
		{
			Key:          "some1",
			LastModified: time.Now().Add(1),
		},
		{
			Key:          "some2",
			LastModified: time.Now().Add(2),
		},
		{
			Key:          "some3",
			LastModified: time.Now().Add(3),
		},
	}
	return objects
}

func TestIsLatest(t *testing.T) {
	objects := generateObjects()

	result := isLatest("some", objects)
	if result == false {
		t.Errorf("Object not found in collection")
	}

	result = isLatest("12312312313", objects)
	if result == true {
		t.Errorf("Object not found in collection")
	}
}
