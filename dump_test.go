package main

import "testing"

func TestExportResult_Filename(t *testing.T) {
	result := &ExportResult{}
	result.Path = "dir/filename.sql"

	filename := result.Filename()

	if filename != "filename.sql" {
		t.Errorf("Filename from path didn't found")
	}
}
