package mysql

import (
	"testing"
)

func TestA(t *testing.T) {
	config := NewDefaultConfig()
	config.Addr = "127.0.0.1:3306"
	config.User = "root"
	config.Passwd = "Bangdao01"
	d, err := NewDB(config)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(d.Exec("CREATE DATABASE IF NOT EXISTS aaa;"))
	}
}
