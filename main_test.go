package sql

import (
	"github.com/aliworkshop/configer"
	"github.com/aliworkshop/dbcore"
	"os"
	"testing"
)

var db dbcore.RDBMS

func TestMain(m *testing.M) {
	registry := configer.New()
	registry.SetConfigType("yaml")
	f, err := os.Open("./config.sample.yaml")
	if err != nil {
		panic("cannot read config: " + err.Error())
	}
	err = registry.ReadConfig(f)
	if err != nil {
		panic("cannot read config" + err.Error())
	}

	db = NewRepository(registry, nil)
	err = db.Initialize()
	if err != nil {
		panic("cannot initialize database: " + err.Error())
	}

	os.Exit(m.Run())
}
