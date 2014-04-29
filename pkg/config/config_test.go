package config

import "testing"

func TestReadCfg(t *testing.T) {
	err := Load("../../conf/app.ini.sample")
	if err != nil {
		t.Error(err)
	}
	t.Logf("cfg: %v", Config)
}
