package utils

import (
	"gopkg.in/yaml.v3"
	"os"

	log "github.com/sirupsen/logrus"

	"main/pkg/global"
	"main/pkg/types"
)

func ParseConfig(configPath string) {
	yamlFile, err := os.ReadFile(configPath)

	if err != nil {
		log.Fatalf("Problem with reading config.yaml file: %v\n", err)
	}

	var config types.Settings

	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		log.Fatalf("Problem with unmarshaling YAML: %v\n", err)
	}

	global.Config = &config
}
