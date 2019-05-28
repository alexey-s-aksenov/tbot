package tbot

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type tbotConfig struct {
	Bot struct {
		UUID  string `yaml:"uuid"`
		Admin string `yaml:"admin"`
	}
	Log struct {
		File string `yaml:"file"`
	}
	Deluge struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
	}
	Weather struct {
		Key string `yaml:"key"`
	}
	Users struct {
		File string `yaml:"file"`
	}
	Proxy struct {
		ProxyURL string `yaml:"proxy"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
}

func readConfig() (*tbotConfig, error) {
	file, err := os.Open("config.yaml")
	if err != nil {
		return new(tbotConfig), err
	}
	data, _ := ioutil.ReadAll(file)
	file.Close()

	var config tbotConfig

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return new(tbotConfig), err
	}

	return &config, nil

}
