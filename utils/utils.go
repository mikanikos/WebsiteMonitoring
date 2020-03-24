package utils

import (
	"fmt"
	"github.com/mikanikos/WebsiteMonitoringTool/common"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"runtime"
)

const configDir = "/configs/"

// parse config file given as a parameter and returns the list of website configurations specified
func ParseConfigFile(fileName string) (*common.Config, error) {

	_, configFileName, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	yamlFile, err := ioutil.ReadFile(path.Dir(path.Dir(configFileName)) + configDir + fileName)
	if err != nil {
		return nil, fmt.Errorf("error while reading file given from input: %s", err)
	}

	config := common.NewConfig()

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, fmt.Errorf("error while parsing config yaml file: %s", err)
	}

	return config, nil
}
