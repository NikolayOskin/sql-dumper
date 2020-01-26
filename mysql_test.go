package main

import "testing"

func TestMySQL_Dump(t *testing.T) {
	MysqlDumpCmd = "true"
	mysql := &MySQL{}
	mysql.Password = "some"
	dump := mysql.Dump()

	if dump.MIME != "application/sql" {
		t.Errorf("wrong mime type")
	}
}
