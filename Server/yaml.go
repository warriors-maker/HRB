package Server

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Yaml struct {
	Trusted []string `yaml:"trusted,flow"`
	Faulty []string `yaml:"faulty,flow"`
	Algorithm int `yaml:"algorithm"`
	Source_Byzantine bool `yaml:"source_byzantine"`
	Data_size int `yaml:"data_size"`
	Rounds int `yaml:"rounds"`
	Faulty_Behavior string `yaml:"faulty_behavior"`
}


func decodeYamlFile() Yaml{
	var yamlStruct Yaml
	config, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err)
	}


	err = yaml.Unmarshal(config, &yamlStruct)
	if err != nil {
		fmt.Println(err)
	}

	return yamlStruct
}