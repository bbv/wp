package main

import (
	"./db"
	"./poller"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Tasks []poller.Task `yaml:"urls",inline`
}

func main() {
	fmt.Println("vim-go")
	appConfig, err := poller.ReadAppConfig("app.yml")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(appConfig)
	config, err := readConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config)

	err = db.Init(appConfig.Db)
	if err != nil {
		log.Fatal(err)
	}
	poller := poller.NewPoller(config.Tasks)
	fmt.Println(poller)
	quit := make(chan int)
	<-quit
}

func readConfig(filename string) (Config, error) {
	cf, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = yaml.UnmarshalStrict(cf, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config, nil
}
