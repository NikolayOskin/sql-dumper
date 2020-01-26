package main

import (
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
