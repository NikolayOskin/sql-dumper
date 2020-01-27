package main

import (
	"bytes"
	"fmt"
	"gopkg.in/amz.v3/s3"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

func TestConvertStringToTime(t *testing.T) {
	s := "2006-01-02T15:04:05.000Z"
	result := convertStringToTime(s)
	parsedStr, _ := time.Parse(s, s)

	if result != parsedStr {
		t.Errorf("time conversion didn't match")
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
	s.BucketName = "some-bucket"
	s.AccessKey = "****"
	s.ClientSecret = "***"

	_, err := s.getBucket()
	if err != nil {
		t.Errorf("BucketName is not returned even if its name is valid")
	}
}

func TestS3_Upload(t *testing.T) {
	s := &S3{
		Region:       "some",
		BucketName:   "some",
		AccessKey:    "some",
		ClientSecret: "some",
	}

	testDumps := []*ExportResult{
		{
			Path:  "./tests/dump.sql",
			Error: nil,
		},
		{
			Path:  "./tests/dump_that_not_existed.sql",
			Error: nil,
		},
		{
			Path:  "",
			Error: nil,
		},
		{
			Path:  "./tests/dump.sql",
			Error: fmt.Errorf(""),
		},
	}

	for _, dump := range testDumps {
		if s.Upload(dump) == nil {
			t.Errorf("failed with %v", dump)
		}
	}
}

func TestUploadWithIncorrectBucketName(t *testing.T) {
	s := &S3{}
	s.BucketName = "some@@@bucket"
	s.AccessKey = "****"
	s.ClientSecret = "***"

	dump := &ExportResult{}
	err := s.Upload(dump)
	if err == nil {
		t.Errorf("Incorrect bucket name didn't throw and error")
	}
}

func TestS3_DeleteOldFiles(t *testing.T) {
	var buf bytes.Buffer
	b := &BucketMock{}
	s := &S3{
		Bucket: b,
	}

	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	s.DeleteOldFiles(1)

	if buf.Len() == 0 {
		t.Errorf("nothing was written to log")
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

func TestGetObjects(t *testing.T) {
	list := &s3.ListResp{
		Contents: []s3.Key{
			{
				Key:          "some",
				LastModified: "2006-01-02T15:04:05.000Z",
			},
		},
	}
	S3 := &S3{}
	objects := S3.getObjects(list)
	if objects == nil {
		t.Fail()
	}
}

// HELPERS & MOCKS
type BucketMock struct{}

func (b *BucketMock) List(prefix string, delim string, marker string, max int) (*s3.ListResp, error) {
	l := &s3.ListResp{}
	l.Contents = []s3.Key{
		{
			Key:          "some1",
			LastModified: "2006-01-02T15:04:05.000Z",
		},
		{
			Key:          "some2",
			LastModified: "2006-01-05T15:04:05.000Z",
		},
		{
			Key:          "some3",
			LastModified: "2006-01-09T15:04:05.000Z",
		},
		{
			Key:          "some4",
			LastModified: "2008-01-09T15:04:05.000Z",
		},
	}
	return l, nil
}
func (b *BucketMock) Del(s string) error {
	return fmt.Errorf("failed to delete %s", s)
}
func (b *BucketMock) PutReader(path string, r io.Reader, length int64, contType string, perm s3.ACL) error {
	return nil
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
