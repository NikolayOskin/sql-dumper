package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "./config.json", "Path to json config")
	flag.Parse()

	config, err := getConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	run(config)
}

func run(config *Config) {
	s3 := newS3(config)
	mysql := newMySQL(config)

	ch := make(chan bool)
	go deleteOldDumps(s3, config.DumpsToKeep, ch)

	dump := mysql.Dump()

	if err := s3.Upload(dump); err != nil {
		log.Printf("Problem(s) while uploading file to S3: %v", err)
	}

	if err := dump.Delete(); err != nil {
		log.Println(err)
	}

	<-ch

	fmt.Println("Success!")
}

func newMySQL(config *Config) *MySQL {
	return &MySQL{
		Host:     config.MysqlHost,
		Port:     config.MysqlPort,
		DB:       config.MysqlDb,
		User:     config.MysqlUser,
		Password: config.MysqlPass,
		Options:  nil,
	}
}

func newS3(config *Config) *S3 {
	return &S3{
		Region:       config.AwsRegion,
		BucketName:   config.AwsBucket,
		AccessKey:    config.AwsKey,
		ClientSecret: config.AwsSecret,
	}
}

func getConfig(configFile string) (*Config, error) {
	config := &Config{}
	if err := config.Parse(configFile); err != nil {
		return nil, err
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}
