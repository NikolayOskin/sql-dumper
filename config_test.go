package main

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {
	config := &Config{}

	// Validation empty config
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
		t.Errorf("Config didn't pass validation even if its filled")
	}

	config.AwsSecret, config.AwsKey, config.AwsBucket, config.AwsRegion,
		config.MysqlUser, config.MysqlPort, config.MysqlDb, config.MysqlHost = " ", " ", " ", " ", " ", " ", " ", " "
	err = config.Validate()
	if err == nil {
		t.Errorf("Trimmed config fields are empty strings")
	}
}

func TestConfigInitWithExistedConfig(t *testing.T) {
	configFile := "./config.example.json"
	config := &Config{}
	err := config.Parse(configFile)
	if err != nil {
		t.Errorf("Not existed config file should return err, nil returned")
	}
}

func TestConfigInitWithNonExistedConfig(t *testing.T) {
	configFile := "./config.example123.json"
	config := &Config{}
	err := config.Parse(configFile)
	if err == nil {
		t.Errorf("Not existed config file should return err, nil returned")
	}
}

func TestIsStringEmpty(t *testing.T) {
	str := "    "
	if isEmpty(str) == false {
		t.Errorf("Empty string is not recognized")
	}
}
