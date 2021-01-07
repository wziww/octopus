package permission

import (
	"octopus/config"
	"octopus/log"
	"octopus/message"
	"os"
	"testing"
)

func TestPermission(t *testing.T) {
	os.Setenv("CONFIG_FILE", "../conf/test.conf")
	config.Init()
	log.Init()
	message.Init()
	Init()
	if len(userGroup) != 2 ||
		userGroup[0].Username != "root" ||
		userGroup[0].Password != "root" ||
		userGroup[0].Permission != PERMISSIONDEV|PERMISSIONMONIT|PERMISSIONEXEC ||
		userGroup[1].Username != "viewer" ||
		userGroup[1].Password != "viewer" ||
		userGroup[1].Permission != PERMISSIONMONIT {
		t.Fatal("userGroup Init test error")
	}
	if Get("4623242325232423706c616365686f6c646572303334333423242325232423262a462a4a297c726f6f74") != userGroup[0] {
		t.Fatal("user Get error")
	}
	if Get("aaa") != nil {
		t.Fatal("user Get error")
	}
	if token, per := Login("root", "root"); token != userGroup[0].Token || per != PERMISSIONDEV|PERMISSIONMONIT|PERMISSIONEXEC {
		t.Fatal("user login error")
	}
	if token, per := Login("root1", "root1"); token != "" || per != PERMISSIONNONE {
		t.Fatal("user login error")
	}
}
