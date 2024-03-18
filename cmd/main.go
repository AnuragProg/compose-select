package main

import (
	"fmt"
	parser "github.com/AnuragProg/compose-select/internal/parser"
)


func main(){
	composeFile, err := parser.NewComposeFile("sample.yaml")
	fmt.Println(composeFile)
	if err != nil{
		panic(err.Error())
	}
	yaml, err := composeFile.GetDependentServicesYAML("app2")
	if err != nil{
		panic(err.Error())
	}
	finalYaml := map[string]interface{}{
		"version": 3,
		"services": yaml,
	}
	composeFile.WriteYAML("output.yaml", finalYaml)
}
