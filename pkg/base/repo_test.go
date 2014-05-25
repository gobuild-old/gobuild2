package base

import "testing"

func TestParseCvsURI(t *testing.T) {
	ci, err := ParseCvsURI("github.com/codeskyblue/fswatch")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ci)
	ci, err = ParseCvsURI("github.com/go-xorm/cmd/xorm")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ci)

	ci, err = ParseCvsURI("code.google.com/p/gcfg")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ci)
}
