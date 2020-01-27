package main

import (
	"testing"
)

//func TestGetConfig(t *testing.T) {
//	if os.Getenv("BE_CRASHER") == "1" {
//		getConfig("./config.json")
//		return
//	}
//	cmd := exec.Command(os.Args[0], "-test.run=TestGetConfig")
//	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
//	err := cmd.Run()
//	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
//		return
//	}
//	t.Fatalf("process ran with err %v, want exit status 1", err)
//}

func TestGetConfig(t *testing.T) {
	// Not existed config
	_, err := getConfig("./notexist.json")

	if err == nil {
		t.Fail()
	}

	// Invalid config
	_, err = getConfig("./tests/config.not_filled.json")

	if err == nil {
		t.Fail()
	}

	// Valid config
	_, err = getConfig("./config.example.json")

	if err != nil {
		t.Fail()
	}
}

func Test_New_MySQL(t *testing.T) {
	config := &Config{
		MysqlHost: "",
		MysqlPort: "",
		MysqlDb:   "",
		MysqlUser: "",
	}
	mysql := newMySQL(config)

	testCases := []struct {
		Field, Actual, Expected string
	}{
		{"db-name", mysql.DB, config.MysqlDb},
		{"db-host", mysql.Host, config.MysqlHost},
		{"db-user", mysql.User, config.MysqlUser},
		{"db-port", mysql.Port, config.MysqlPort},
	}
	for _, tc := range testCases {
		if tc.Actual != tc.Expected {
			t.Errorf("failed mySQL struct field %s: actual %s, expected %s", tc.Field, tc.Actual, tc.Expected)
		}
	}
}

func Test_New_S3(t *testing.T) {
	config := &Config{
		AwsRegion: "",
		AwsBucket: "",
		AwsKey:    "",
		AwsSecret: "",
	}
	s3 := newS3(config)

	testCases := []struct {
		Field, Actual, Expected string
	}{
		{"aws-region", s3.Region, config.AwsRegion},
		{"aws-bucket", s3.BucketName, config.AwsBucket},
		{"aws-key", s3.AccessKey, config.AwsKey},
		{"aws-secret", s3.ClientSecret, config.AwsSecret},
	}
	for _, tc := range testCases {
		if tc.Actual != tc.Expected {
			t.Errorf("failed S3 struct field %s: actual %s, expected %s", tc.Field, tc.Actual, tc.Expected)
		}
	}
}
