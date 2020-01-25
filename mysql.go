package main

import (
	"fmt"
	"os/exec"
	"time"
)

var (
	// MysqlDumpCmd is the path to the `mysqldump` executable
	MysqlDumpCmd = "mysqldump"
)

// MySQL is an `Exporter` interface that backs up a MySQL database via the `mysqldump` command
type MySQL struct {
	// DB Host (e.g. 127.0.0.1)
	Host string
	// DB Port (e.g. 3306)
	Port string
	// DB Name
	DB string
	// DB User
	User string
	// DB Password
	Password string
	// Extra mysqldump options
	// e.g []string{"--extended-insert"}
	Options []string
}

// Dump produces a `mysqldump` of the specified database.
func (x MySQL) Dump() *ExportResult {
	result := &ExportResult{MIME: "application/sql"}

	result.Path = fmt.Sprintf(`%v_%v.sql`, time.Now().Format("2006-01-02"), time.Now().Unix())

	options := append(x.dumpOptions(), fmt.Sprintf(`-r%v`, result.Path))

	_, err := exec.Command(MysqlDumpCmd, options...).Output()
	if err != nil {
		err = fmt.Errorf("trying execute mysqldump command: %v", err)
		result.Error = err
	}

	return result
}

func (x MySQL) dumpOptions() []string {
	options := x.Options
	options = append(options, fmt.Sprintf(`-h%v`, x.Host))
	options = append(options, fmt.Sprintf(`-P%v`, x.Port))
	options = append(options, fmt.Sprintf(`-u%v`, x.User))

	if x.Password != "" {
		options = append(options, fmt.Sprintf(`-p%v`, x.Password))
	}
	options = append(options, x.DB)

	return options
}
