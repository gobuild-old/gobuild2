package base

import (
	"strings"
	"testing"
)
import "github.com/remogatto/prettytest"

// Start of setup
type testSuite struct {
	prettytest.Suite
}

func TestRunner(t *testing.T) {
	prettytest.Run(
		t,
		new(testSuite),
	)
}

var v = map[string]string{
	"Name": "codeskyblue",
}

var expect = `{"Name":"codeskyblue"}`

func (t *testSuite) TestObjc2Str() {
	s := Objc2Str(v)
	t.Equal(expect, s)
	t.True(strings.Contains(s, "codeskyblue"))
}

func (t *testSuite) TestObjc2Json() {
	var o map[string]string
	err := Str2Objc(expect, &o)
	t.Nil(err)
	t.Equal(v["Name"], o["Name"])
}
