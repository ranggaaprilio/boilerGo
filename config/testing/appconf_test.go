package testing

import (
	"testing"

	"github.com/ranggaaprilio/boilerGo/config"
)

func TestLoadconf(t *testing.T) {
	conf := config.Loadconf()
	t.Log(conf)
}

func TestDbCon(t *testing.T) {
	config.DbInit()
	conf := config.CreateCon()
	t.Log(conf)
}
