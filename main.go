package main

import (
	"flag"
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

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "./config.json", "Path to json config file")
	flag.Parse()

	config := &Config{}
	config.Init(configFile)

	s3 := &S3{
		Region:       config.AwsRegion,
		Bucket:       config.AwsBucket,
		AccessKey:    config.AwsKey,
		ClientSecret: config.AwsSecret,
	}
	mysql := &MySQL{
		Host:     config.MysqlHost,
		Port:     config.MysqlPort,
		DB:       config.MysqlDb,
		User:     config.MysqlUser,
		Password: config.MysqlPass,
		Options:  nil,
	}

	go deleteOldDumps(s3, config.dumpNames)

	dump := mysql.Dump(config.DumpNameFormat)
	err := s3.Upload(dump)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Deleting dump
	err = os.Remove(dump.Path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Success!")
}

func deleteOldDumps(s3 *S3, dumpnames []string) {
	err := s3.DeleteFilesExcept(dumpnames)
	if err != nil {
		fmt.Println(err.Error())
	}
}
