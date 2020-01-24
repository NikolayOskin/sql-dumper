package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	DumpsToKeep uint `json:"dumps_to_keep"`

	// By default "2006-01-02" format will store dumps like "2020-01-01.sql"
	DumpNameFormat string `json:"dump_name_format"`
}

func (c *Config) Init(configFile string) error {

	configJson, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("reading json config file: %v", err)
	}

	err = json.Unmarshal(configJson, &c)
	if err != nil {
		return fmt.Errorf("trying to unmarshal json config: %v", err)
	}

	return nil
}
