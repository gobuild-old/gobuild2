package base

import (
	"encoding/json"

	"github.com/gobuild/log"
)

func Str2Objc(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

func Objc2Str(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Errorf("json2str got error: %v", err)
		return ""
	}
	return string(data)
}
