package config

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	configurationData	*conf
)

func init() {
	configurationData = &conf{}
}

func parseAndOverwriteConfig(yamlFile []byte) error {
	newConfigurationData := &conf{}
	err := yaml.UnmarshalStrict(yamlFile, newConfigurationData)
	if err != nil {
		return err
	}

	//TODO: set default values for unconfigured fields
	//TODO: validate config values

	configurationData = newConfigurationData
	return nil
}

func configuration() *conf {
	return configurationData
}

func LoadConfiguration(configurationFilename string) error {
	if configurationFilename == "" {
		return errors.New("no configuration file name provided")
	}

	log.Printf("[00000000] INFO Reading configuration at %s ...", configurationFilename)
	yamlFile, err := ioutil.ReadFile(configurationFilename)
	if err != nil {
		return err
	}

	err = parseAndOverwriteConfig(yamlFile)
	return err
}