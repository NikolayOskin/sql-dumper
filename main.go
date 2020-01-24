package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		configFile string
		err        error
	)

	flag.StringVar(&configFile, "config", "./config.json", "Path to json config")
	flag.Parse()

	config := &Config{}
	err = config.Init(configFile)
	if err != nil {
		fmt.Printf("Error while initializing configuration, exiting: %v", err)
		os.Exit(1)
	}

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

	ch := make(chan bool)
	go deleteOldDumps(s3, config.DumpsToKeep, ch)

	dump := mysql.Dump(config.DumpNameFormat)

	err = s3.Upload(dump)
	if err != nil {
		fmt.Printf("Problem(s) while uploading file to S3: %v", err)
	}

	err = dump.Delete()
	if err != nil {
		fmt.Println(err)
	}

	<-ch

	fmt.Println("Success!")
}

func deleteOldDumps(s3 *S3, dumpsToKeep uint, ch chan<- bool) {
	err := s3.DeleteFilesExceptLatest(dumpsToKeep)
	if err != nil {
		fmt.Printf("Problem(s) with deleting old dumps from S3: %v", err)
	}

	ch <- true
}
