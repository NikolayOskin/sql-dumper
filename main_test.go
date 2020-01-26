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
