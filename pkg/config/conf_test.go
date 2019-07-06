package config

import (
	"fmt"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestDefaults(t *testing.T) {
	cas := load(t, "cassandra")
	def := load(t, "default")
	if cas != def {
		t.Errorf("default: %s\n\ncassandra: %s\n\n", def, cas)
	}
}

func load(t *testing.T, name string) string {
	conf := CassandraYAML{}
	b, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.yaml", name))
	if err != nil {
		t.Fatalf("Error opening %s.yaml: %v", name, err)
	}
	if err := yaml.Unmarshal(b, &conf); err != nil {
		t.Fatalf("Error in unmarshal: %v", err)
	}
	yb, err := yaml.Marshal(&conf)
	if err != nil {
		t.Fatalf("Error in marshal: %v", err)
	}
	return string(yb)
}
