package main

import (
	"fmt"
	"log"
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
		return fmt.Errorf("failed to delete dump file %s: %s", x.Path, err)
	}
	return nil
}

func deleteOldDumps(s3 *S3, dumpsToKeep int, ch chan<- bool) {
	if dumpsToKeep > 0 {
		err := s3.DeleteOldFiles(dumpsToKeep)
		if err != nil {
			log.Printf("failed to delete old dumps from S3: %v", err)
		}
	}

	ch <- true
}
