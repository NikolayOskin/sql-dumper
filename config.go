package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
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
}

func (c *Config) Parse(configFile string) error {

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

func (c *Config) Validate() error {
	if isEmpty(c.MysqlHost) || isEmpty(c.MysqlPort) || isEmpty(c.MysqlDb) || isEmpty(c.MysqlUser) {
		return fmt.Errorf("MySQL data is not filled")
	}
	if isEmpty(c.AwsRegion) || isEmpty(c.AwsBucket) || isEmpty(c.AwsKey) || isEmpty(c.AwsSecret) {
		return fmt.Errorf("AWS data is not filled")
	}

	return nil
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) < 1
}
