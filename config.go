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

	DumpsToKeep int `json:"dumps_to_keep"`
}

func (c *Config) Parse(configFile string) error {

	configJson, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed parsing json config file %s: %s", configFile, err)
	}

	err = json.Unmarshal(configJson, &c)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json config %s: %s", configFile, err)
	}

	return nil
}

func (c *Config) Validate() error {
	if isEmpty(c.MysqlHost) {
		return fmt.Errorf("db host host is not filled")
	}
	if isEmpty(c.MysqlPort) {
		return fmt.Errorf("db port is not filled")
	}
	if isEmpty(c.MysqlDb) {
		return fmt.Errorf("db name is not filled")
	}
	if isEmpty(c.MysqlUser) {
		return fmt.Errorf("db user is not filled")
	}
	if isEmpty(c.AwsRegion) || isEmpty(c.AwsBucket) || isEmpty(c.AwsKey) || isEmpty(c.AwsSecret) {
		return fmt.Errorf("AWS data is not filled")
	}

	if c.DumpsToKeep < 0 {
		return fmt.Errorf("dumpsToKeep should be >= 0, %v given", c.DumpsToKeep)
	}

	return nil
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) < 1
}
