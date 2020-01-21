package main

import (
	"flag"
	"fmt"
)

func main() {
	var configFile string

	flag.StringVar(&configFile, "config", "./config.json", "Path to json config")
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

	ch := make(chan bool)
	go deleteOldDumps(s3, config.dumpNames, ch)

	dump := mysql.Dump(config.DumpNameFormat)

	err := s3.Upload(dump)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = dump.Delete()
	if err != nil {
		fmt.Println(err)
	}

	<-ch

	fmt.Println("Success!")
}

func deleteOldDumps(s3 *S3, latestDumps []string, ch chan<- bool) {
	err := s3.DeleteFilesExcept(latestDumps)
	if err != nil {
		fmt.Println(err.Error())
	}

	ch <- true
}
