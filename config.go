package gogomailer

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

const (
	CONFIGPATH = "config.yaml"
)

var (
	Conf map[string]map[string]string
)

func loadConfig() (map[string]map[string]string, error) {
	var config map[string]map[string]string

	configBytes, err := ioutil.ReadFile(CONFIGPATH)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load the config: %s", CONFIGPATH))
		panic(err)
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to unmarshal the config: %s", CONFIGPATH))
		return config, err
	}

	log.Println(fmt.Sprintf("Config loaded successfully: %s", CONFIGPATH))
	return config, err
}
