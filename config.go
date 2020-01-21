package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Config struct {
	MysqlHost string `json:"mysql_host"`
	MysqlPort string `json:"mysql_port"`
	MysqlDb   string `json:"mysql_db"`
	MysqlUser string `json:"mysql_user"`
	MysqlPass string `json:"mysql_password"`

	AwsRegion string `json:"aws_region"`
	AwsBucket string `json:"aws_bucket"`
	AwsKey    string `json:"aws_key"`
	AwsSecret string `json:"aws_secret"`

	DumpsToKeep int `json:"dumps_to_keep"`

	// By default "2006-01-02" format will store dumps like "2020-01-01.sql"
	DumpNameFormat string `json:"dump_name_format"`
	dumpNames      []string
}

func (c *Config) Init(configFile string) error {
	var date, dumpName string

	configJson, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("reading json config file: %v", err)
	}

	err = json.Unmarshal(configJson, &c)
	if err != nil {
		return fmt.Errorf("trying to unmarshal json config: %v", err)
	}

	// collecting filenames which will be kept on S3 bucket
	for i := 1; i <= c.DumpsToKeep; i++ {
		date = time.Now().Add(time.Hour * time.Duration(i*-24)).Format(c.DumpNameFormat)
		dumpName = fmt.Sprintf(`%v.sql`, date)
		c.dumpNames = append(c.dumpNames, dumpName)
	}

	return nil
}
