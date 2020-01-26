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

	dump := mysql.Dump()

	if err := s3.Upload(dump); err != nil {
		fmt.Printf("Problem(s) while uploading file to S3: %v", err)
	}

	if err := dump.Delete(); err != nil {
		fmt.Println(err)
	}

	<-ch

	fmt.Println("Success!")
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
