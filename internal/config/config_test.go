package config

import "testing"

func TestConfigCreate(t *testing.T) {
	k1 := "debug"
	v1 := false
	k2 := "dbParams"
	v2 := map[string]string{
		"login":    "admin",
		"password": "testing",
		"port":     "12345",
	}
	k3 := "dbTable"
	v3 := map[string]map[string]string{
		"id": map[string]string{
			"type":       "int",
			"primaryKey": "yes",
		},
		"shortLink": map[string]string{
			"type":   "text",
			"unique": "true",
		},
		"origLink": map[string]string{
			"type":   "text",
			"unique": "true",
		},
	}
	config := Config{
		k1: v1,
		k2: v2,
		k3: v3,
	}
	if config[k1].(bool) {
		t.Error("expected false but got true on k1")
	}
	if config[k2].(map[string]string)["login"] != "admin" {
		t.Error("unexpected string from map[string]string")
	}
	if config[k3].(map[string]map[string]string)["shortLink"]["type"] != "text" {
		t.Error("unexpected string from map[string]map[string]string")
	}
}
