package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// ExportResult the result of dump db operation
type ExportResult struct {
	// Path to exported file
	Path string
	// MIME type of the exported file (e.g. application/x-tar)
	MIME  string
	Error error
}

// Filename returns the just filename component of the `Path` attribute
func (x ExportResult) Filename() string {
	_, filename := filepath.Split(x.Path)
	return filename
}

// Delete dump from server
func (x *ExportResult) Delete() error {
	err := os.Remove(x.Path)
	if err != nil {
		return fmt.Errorf("trying to delete dump file from server: %v", err)
	}
	return nil
}

func deleteOldDumps(s3 *S3, dumpsToKeep int, ch chan<- bool) {
	err := s3.DeleteOldFiles(dumpsToKeep)
	if err != nil {
		fmt.Printf("Problem(s) with deleting old dumps from S3: %v", err)
	}

	ch <- true
}
