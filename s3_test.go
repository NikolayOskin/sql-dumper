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

func Test_Convert_LastModified_String_To_Time(t *testing.T) {
	s := "2006-01-02T15:04:05.000Z"
	result := convertStringToTime(s)
	parsedStr, _ := time.Parse(s, s)

	if result != parsedStr {
		t.Errorf("time conversion didn't match")
	}
}

func Test_Sort_Files_By_LastModified(t *testing.T) {
	// Prepare
	obj1 := ObjectS3{
		Key:          "some_name1",
		LastModified: time.Now(),
	}
	obj2 := ObjectS3{
		Key:          "some_name2",
		LastModified: time.Now().Add(100),
	}
	objects := []ObjectS3{obj1, obj2}

	// Sort
	firstBefore := objects[0]

	sortByLastModified(objects)

	firstAfter := objects[0]

	// Check
	if firstBefore == firstAfter {
		t.Errorf("Objects in slice were not sorted")
	}
}

func Test_GetBucket(t *testing.T) {
	s := &S3{}
	s.BucketName = "some-bucket"
	s.AccessKey = "****"
	s.ClientSecret = "***"

	_, err := s.getBucket()
	if err != nil {
		t.Errorf("BucketName is not returned even if its name is valid")
	}
}

func Test_Upload_Returns_Errors(t *testing.T) {
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

func Test_Upload_With_Incorrect_BucketName(t *testing.T) {
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

func Test_DeleteOldFiles_Writes_Deleting_Errors_To_Log(t *testing.T) {
	var buf bytes.Buffer
	b := &BucketMock{}
	s := &S3{Bucket: b}

	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	s.DeleteOldFiles(1)

	if buf.Len() == 0 {
		t.Errorf("nothing was written to log")
	}
}

func Test_KeepLatest(t *testing.T) {
	objects := generateObjects() // 4 objects

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

func Test_IsLatest(t *testing.T) {
	objects := generateObjects()

	testCases := []struct {
		Name     string
		Expected bool
	}{
		{"some", true},
		{"blahblah", false},
	}

	for _, tc := range testCases {
		result := isLatest(tc.Name, objects)
		if result != tc.Expected {
			t.Errorf("failed, name:%s, expected: %t, actual: %t", tc.Name, tc.Expected, result)
		}
	}
}

func Test_GetObjects_From_Bucket(t *testing.T) {
	// Mocking
	list := &s3.ListResp{
		Contents: []s3.Key{
			{
				Key:          "some",
				LastModified: "2006-01-02T15:04:05.000Z",
			},
		},
	}
	S3 := &S3{}

	// Trying to retrieve objects from mocked response
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
