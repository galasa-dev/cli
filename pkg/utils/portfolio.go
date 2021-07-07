//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package utils

import (
	"bufio"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Portfolio struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   PortfolioMetadata `yaml:"metadata"`

	Classes    []PortfolioClass  `yaml:"classes"`    
}

type PortfolioMetadata struct {
	Name string  `yaml:"name"`
}

type PortfolioClass struct {
	Bundle   string   `yaml:"bundle"`
	Class    string   `yaml:"class"`
	Stream   string   `yaml:"stream"`
}


func CreatePortfolio(testSelection *TestSelection) Portfolio {
	portfolio := Portfolio {
		APIVersion: "v1alpha",
		Kind:       "galasa.dev/testPortfolio",
		Metadata: PortfolioMetadata { Name: "adhoc"},
	}

	portfolio.Classes = make([]PortfolioClass, 0)

	for _, selectedClass := range testSelection.Classes {
		portfolioClass := PortfolioClass {
			Bundle: selectedClass.Bundle,
			Class: selectedClass.Class,
			Stream: selectedClass.Stream,
		}
		portfolio.Classes = append(portfolio.Classes, portfolioClass)
	}

	return portfolio
}


func WritePortfolio(portfolio Portfolio, filename string) {

	file,err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(file)

	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)

	err = encoder.Encode(&portfolio)
	if err != nil {
		panic(err)
	}
	w.Flush()
	encoder.Close()
	file.Close()
}


func LoadPortfolio(filename string) Portfolio {

	var portfolio Portfolio

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(b, &portfolio)
	if err != nil {
		panic(err)
	}

	if portfolio.APIVersion != "v1alpha" {
		panic("Portfolio file is not version 'v1alpha'")
	}

	if portfolio.Kind != "galasa.dev/testPortfolio" {
		panic("Portfolio file is not kind 'galasa.dev/testPortfolio'")
	}

	return portfolio
}
