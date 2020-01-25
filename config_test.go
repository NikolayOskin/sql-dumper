package main

import (
	"testing"
)

func TestConfigValidation(t *testing.T) {

	testCases := []*Config{
		{},
		{
			MysqlHost: "",
			MysqlDb:   "",
			MysqlPort: "",
			MysqlUser: "",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "  ",
			MysqlDb:   "  ",
			MysqlPort: "  ",
			MysqlUser: "  ",
			MysqlPass: "  ",
			AwsRegion: "  ",
			AwsBucket: "  ",
			AwsKey:    "  ",
			AwsSecret: "  ",
		},
		{
			MysqlHost: "",
			MysqlDb:   "",
			MysqlPort: "",
			MysqlUser: "",
			MysqlPass: "",
		},
		{
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "",
			MysqlPort: "",
			MysqlUser: "",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "",
			MysqlDb:   "name",
			MysqlPort: "",
			MysqlUser: "",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "",
			MysqlDb:   "",
			MysqlPort: "3232",
			MysqlUser: "",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "",
			MysqlDb:   "",
			MysqlPort: "",
			MysqlUser: "eqwe",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "name",
			MysqlPort: "",
			MysqlUser: "",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "",
			MysqlPort: "3234",
			MysqlUser: "",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "",
			MysqlPort: "",
			MysqlUser: "ewer",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "name",
			MysqlPort: "3306",
			MysqlUser: "",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "name",
			MysqlPort: "",
			MysqlUser: "some",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "name",
			MysqlPort: "3306",
			MysqlUser: "some",
			MysqlPass: "",
			AwsRegion: "123",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "name",
			MysqlPort: "3306",
			MysqlUser: "some",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "123",
			AwsKey:    "",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "name",
			MysqlPort: "3306",
			MysqlUser: "some",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "123",
			AwsSecret: "",
		},
		{
			MysqlHost: "localhost",
			MysqlDb:   "name",
			MysqlPort: "3306",
			MysqlUser: "some",
			MysqlPass: "",
			AwsRegion: "",
			AwsBucket: "",
			AwsKey:    "",
			AwsSecret: "123",
		},
	}

	for _, config := range testCases {
		err := config.Validate()
		if err == nil {
			t.Errorf("Validation rule is not correct, %v", config)
		}
	}
	config := &Config{
		MysqlHost:   "localhost",
		MysqlPort:   "3306",
		MysqlDb:     "some",
		MysqlUser:   "some",
		MysqlPass:   "",
		AwsRegion:   "some",
		AwsBucket:   "some",
		AwsKey:      "some",
		AwsSecret:   "som",
		DumpsToKeep: 5,
	}
	err := config.Validate()
	if err != nil {
		t.Errorf("Valid config throwed error")
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
