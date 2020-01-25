package main

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	config := &Config{}

	err := config.Validate()
	if err == nil {
		t.Errorf("Validation doesn't work")
	}

	config.MysqlHost, config.MysqlUser, config.MysqlDb, config.MysqlPort = "s", "s", "s", "s"

	err = config.Validate()
	if err == nil {
		t.Errorf("Validation doesn't work")
	}

	config.AwsSecret, config.AwsKey, config.AwsBucket, config.AwsRegion,
		config.MysqlUser, config.MysqlPort, config.MysqlDb, config.MysqlHost = "s", "s", "s", "s", "s", "s", "s", "s"
	err = config.Validate()
	if err != nil {
		t.Errorf("Config should be validated if all required fields are filled")
	}
}

func TestConfigInitWithExistedConfig(t *testing.T) {
	configFile := "./config.example.json"
	config := &Config{}
	err := config.Init(configFile)
	if err != nil {
		t.Errorf("Not existed config file should return err, nil returned")
	}
}

func TestConfigInitWithNonExistedConfig(t *testing.T) {
	configFile := "./config.example123.json"
	config := &Config{}
	err := config.Init(configFile)
	if err == nil {
		t.Errorf("Not existed config file should return err, nil returned")
	}
}
