package main

import (
	"testing"
)

func TestConfigInit(t *testing.T) {
	configFile := "./config.json"
	config := &Config{}
	config.Init(configFile)

	if len(config.dumpNames) != config.DumpsToKeep {
		t.Errorf("Dumpnames slice should be length %d, got: %d.", config.DumpsToKeep, len(config.dumpNames))
	}

}
