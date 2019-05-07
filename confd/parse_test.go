package confd

import "testing"

func TestParseYAMLFile(t *testing.T) {
	cfg, err := parseYamlFile("../example.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(cfg)
}
