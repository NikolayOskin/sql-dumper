package main

import (
	"fmt"
	"time"
)

type Config struct {
	// how many dump files will be kept on S3 bucket by default
	dumpsToKeep int
	// filename date format, example "2006-05-15"
	fileNameFormat string
	dumpNames      []string
}

func (c *Config) init() {
	var date, fileName string

	// collecting filenames which will be kept on S3 bucket
	for i := 1; i <= c.dumpsToKeep; i++ {
		date = time.Now().Add(time.Hour * time.Duration(i*-24)).Format(c.fileNameFormat)
		fileName = fmt.Sprintf(`%v.sql`, date)
		c.dumpNames = append(c.dumpNames, fileName)
	}
}
