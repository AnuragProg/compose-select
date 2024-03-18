package parser

import (
	"errors"
	"os"
	"io/fs"

	"gopkg.in/yaml.v3"
)

type ComposeFile struct {
	services map[string]interface{}
}

func NewComposeFile(filename string) (*ComposeFile, error){
	file, err := os.Open(filename)
	if err != nil{
		return nil, err
	}

	var composeYAML map[string]interface{}
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&composeYAML); err != nil{
		return nil, err
	}

	services, ok := composeYAML["services"].(map[string]interface{})
	if !ok {
		return nil, errors.New("services section not defined in compose file")
	}

	return &ComposeFile{
		services: services,
	}, nil
}

func (cf *ComposeFile) GetServiceNames() []string {
	var serviceNames []string
	for service := range cf.services {
		serviceNames = append(serviceNames, service)
	}
	return serviceNames
}

func (cf *ComposeFile) GetDependencyNamesForService(serviceName string) ([]string, error) {
	serviceIface, ok := cf.services[serviceName]
	if !ok {
		return nil, errors.New("service not found")
	}
	var dependencyNames []string
	if serviceIface == nil{
		return dependencyNames, nil
	}

	service, _ := serviceIface.(map[string]interface{})
	dependencies, _ := service["depends_on"].([]interface{})
	for _, dep := range dependencies {
		if depStr, ok := dep.(string); ok {
			dependencyNames = append(dependencyNames, depStr)
		}
	}
	return dependencyNames, nil
}

func (cf *ComposeFile) GetDependentServicesYAML(serviceName string) (map[string]interface{}, error) {
	yaml := map[string]interface{}{}

	serviceNames, err := cf.GetDependencyNamesForService(serviceName)
	if err != nil{
		return nil, err
	}

	for len(serviceNames) > 0 {

		// pop from services
		curServiceName := serviceNames[0]
		serviceNames = serviceNames[1:]

		// construct service yaml
		curServiceYAML, ok := cf.services[curServiceName]
		if !ok {
			return yaml, errors.New(curServiceName + " not found")
		}
		yaml[curServiceName] = curServiceYAML // add service yaml

		// retrieve dependencies
		curServiceDependencies, err := cf.GetDependencyNamesForService(curServiceName)
		if err != nil{
			return nil, err
		}

		// add non visited dependency
		for _, nextService := range curServiceDependencies{
			if _, ok := yaml[nextService]; ok || nextService == serviceName{
				continue
			}
			serviceNames = append(serviceNames, nextService)
		}
	}
	return yaml, nil
}


func (cf *ComposeFile) WriteYAML(filename string, yamlData map[string]interface{}) error {

	// marshal yaml data
	yamlSerialData, err := yaml.Marshal(yamlData)
	if err != nil{
		return err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.ModePerm)	
	if err != nil{
		return err
	}
	defer file.Close()


	_, err = file.Write(yamlSerialData)
	if err != nil{
		return err
	}

	return nil
}







