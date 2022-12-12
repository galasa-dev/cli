/*
 * Copyright contributors to the Galasa project
 */
package utils

import (
	"bufio"
	"io/ioutil"
	"os"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Portfolio struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   PortfolioMetadata `yaml:"metadata"`

	Classes []PortfolioClass `yaml:"classes"`
}

type PortfolioMetadata struct {
	Name string `yaml:"name"`
}

type PortfolioClass struct {
	Bundle    string            `yaml:"bundle"`
	Class     string            `yaml:"class"`
	Stream    string            `yaml:"stream"`
	Overrides map[string]string `yaml:"overrides"`
}

func NewPortfolio() Portfolio {
	portfolio := Portfolio{
		APIVersion: "v1alpha",
		Kind:       "galasa.dev/testPortfolio",
		Metadata:   PortfolioMetadata{Name: "adhoc"},
	}

	portfolio.Classes = make([]PortfolioClass, 0)

	return portfolio
}

// Inside the portfolio file, it must carry a format field with this value inside.
const PORTOLIO_DECLARED_FORMAT_VERSION = "v1alpha"

// Inside the portfolio file, it should claim to be a resource of this kind.
const PORTFOLIO_DECLARED_RESOURCE_KIND = "galasa.dev/testPortfolio"

func CreatePortfolio(testSelection *TestSelection, testOverrides *map[string]string, portfolio *Portfolio) {

	for _, selectedClass := range testSelection.Classes {
		portfolioClass := PortfolioClass{
			Bundle:    selectedClass.Bundle,
			Class:     selectedClass.Class,
			Stream:    selectedClass.Stream,
			Overrides: *testOverrides,
		}
		portfolio.Classes = append(portfolio.Classes, portfolioClass)
	}
}

func WritePortfolio(portfolio Portfolio, filename string) {

	file, err := os.Create(filename)
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
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_OPEN_PORTFOLIO_FILE_FAILED, filename, err.Error())
		panic(err)
	}

	err = yaml.Unmarshal(b, &portfolio)
	if err != nil {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PORTFOLIO_BAD_FORMAT, filename, err.Error())
		panic(err)
	}

	// Check the portfolio file claims to be the correct format version.
	if portfolio.APIVersion != PORTOLIO_DECLARED_FORMAT_VERSION {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PORTFOLIO_BAD_FORMAT_VERSION, filename, PORTOLIO_DECLARED_FORMAT_VERSION)
		panic(err)
	}

	if portfolio.Kind != PORTFOLIO_DECLARED_RESOURCE_KIND {
		err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_PORTFOLIO_BAD_RESOURCE_KIND, filename, PORTFOLIO_DECLARED_RESOURCE_KIND)
		panic(err)
	}

	return portfolio
}
