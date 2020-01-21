package main

import (
	"testing"
)

func TestConfigInit(t *testing.T) {
	config := &Config{
		dumpsToKeep:    10,
		fileNameFormat: "2006-01-02",
	}
	config.init()
	if len(config.dumpNames) != config.dumpsToKeep {
		t.Errorf("Dumpnames slice should be length %d, got: %d.", config.dumpsToKeep, len(config.dumpNames))
	}
}
