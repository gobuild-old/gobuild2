package main

import "testing"

func TestReadCfg(t *testing.T) {
	cfg, err := readCfg("conf/app.ini.sample")
	if err != nil {
		t.Error(err)
	}
	t.Logf("cfg: %v", cfg)
}
