package main

import "testing"

func TestMySQL_Dump(t *testing.T) {
	MysqlDumpCmd = "true"
	mysql := &MySQL{}
	dump := mysql.Dump()

	if dump.MIME != "application/sql" {
		t.Errorf("wrong mime type")
	}
}
