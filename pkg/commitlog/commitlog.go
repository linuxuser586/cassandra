package commitlog

import (
	"log"

	"github.com/spf13/afero"

	"github.com/magiconair/properties"
)

const (
	propFile       = "/etc/cassandra/commitlog_archiving.properties"
	customPropFile = "/conf/commitlog_archiving.properties"
)

// Update commitlog_archiving.properties
func Update() {
	c := load(customPropFile)
	if c != nil {
		o := load(propFile)
		o.Merge(c)
		save(o)
	}
}

func load(f string) *properties.Properties {
	exists, err := afero.Exists(fs, f)
	if err != nil {
		log.Fatal(err)
	}
	if !exists {
		return nil
	}
	b, err := afero.ReadFile(fs, f)
	if err != nil {
		log.Fatal(err)
	}
	p, err := properties.Load(b, properties.UTF8)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

func save(p *properties.Properties) {
	if err := afero.WriteFile(fs, propFile, []byte(p.String()), 0644); err != nil {
		log.Fatal(err)
	}
}
