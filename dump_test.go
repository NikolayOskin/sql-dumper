package main

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestExportResult_Filename(t *testing.T) {
	result := &ExportResult{}
	result.Path = "dir/filename.sql"

	filename := result.Filename()

	if filename != "filename.sql" {
		t.Errorf("Filename from path didn't found")
	}
}

func TestExportResult_Delete(t *testing.T) {
	file, _ := os.Create("./tests/test.txt")
	file.Close()

	result := &ExportResult{}
	result.Path = "./tests/test.txt"

	err := result.Delete()
	if err != nil {
		t.Fatal(err)
	}

	result.Path = "./tests/not_exist.txt"

	err = result.Delete()
	if err == nil {
		t.Fatal(err)
	}
}

func Test_Delete_Old_Dumps(t *testing.T) {
	var buf bytes.Buffer

	s3 := &S3{}
	ch := make(chan bool)

	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	go deleteOldDumps(s3, 10, ch)

	<-ch
	if buf.Len() == 0 {
		t.Errorf("nothing was written to log")
	}
}
