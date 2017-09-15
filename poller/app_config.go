package poller

import (
	"../db"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type AppConfig struct {
	Port int
	Db   db.DBConfig
}

func ReadAppConfig(filename string) (AppConfig, error) {
	cf, err := ioutil.ReadFile(filename)
	if err != nil {
		return AppConfig{}, err
	}
	config := AppConfig{}
	err = yaml.UnmarshalStrict(cf, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config, nil
}
