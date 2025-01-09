package config

import "testing"

func TestConfig(t *testing.T) {
	config := New("./config.yaml")
	t.Logf("%+v\n", config)
}
