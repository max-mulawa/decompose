package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type composeFile struct {
	Version  string                 `yaml:"version"`
	Services map[string]interface{} `yaml:"services,omitempty"`
}

func main() {
	//command : service
	profileName := "tests"
	includeNotProfiled := true
	composeFilePath := "../../test/devcontainers/example1-compose.yml"
	outputFilePath := "./out.yml"

	//command : devcontainer
	// serviceName (current top folder name)
	//

	file, err := os.ReadFile(composeFilePath)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(file, data)
	if err != nil {
		log.Fatal("failed to unmarshal")
	}

	matchingSrvs := make(map[string]interface{})

	srv, srvExists := data["services"].(map[string]interface{})
	if srvExists {
		for ksrv, vsrv := range srv {
			serviceMap := vsrv.(map[string]interface{})
			profiles, prfExists := serviceMap["profiles"].([]interface{})
			if prfExists {
				for _, prf := range profiles {
					if prf == profileName {
						fmt.Printf("service %s is matching profile\n", ksrv)
						delete(serviceMap, "profiles")
						matchingSrvs[ksrv] = serviceMap
					}
				}
			} else if includeNotProfiled {
				fmt.Printf("service %s with no profiles is matching\n", ksrv)
				matchingSrvs[ksrv] = serviceMap
			}

		}
	}

	output := &composeFile{
		Version:  "3",
		Services: matchingSrvs,
	}
	outputYaml, err := yaml.Marshal(output)
	if err != nil {
		log.Fatalf("failed to marshal yaml: %v", err)
	}

	err = os.WriteFile(outputFilePath, outputYaml, 0664)
	if err != nil {
		log.Fatalf("failed to persist yaml: %v", err)
	}
}
