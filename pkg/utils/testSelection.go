//
// Licensed Materials - Property of IBM
//
// (c) Copyright IBM Corp. 2021.
//

package utils

import (
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	packages       *[]string
	regexSelect    *bool
)

type TestSelection struct {
	Classes []TestClass
}

type TestClass struct {
	Bundle string
	Class  string
	Stream string
}


func AddCommandFlags(command *cobra.Command) {
	packages = command.Flags().StringSlice("packages", make([]string, 0), "packages of which tests will be selected from, packages are selected if the name contains this string, or if --regex is specified then matches the regex")
	regexSelect = command.Flags().Bool("regex", false, "Test selection is performed by using regex")
}


func SelectTests(testCatalog TestCatalog, stream string) TestSelection {

	testSelection := TestSelection{ Classes: make([]TestClass, 0)}

	selectTestsByPackage(testCatalog, &testSelection, stream)

	return testSelection
}


func selectTestsByPackage(testCatalog TestCatalog, testSelection *TestSelection, stream string) {

	if len(*packages) < 1 {
		return 
	}

	if testCatalog["packages"] == nil {
		return
	}

	regexPatterns := make([]*regexp.Regexp, 0)
	if *regexSelect {
		// Create patterns of actual regex 
		for _, selectionPackage := range *packages {
			r, err := regexp.Compile(selectionPackage)
			if err != nil {
				fmt.Printf("Error with regex '%v' - %v", selectionPackage, err)
				panic(err)
			}
			regexPatterns = append(regexPatterns, r)
		}
	} else {
		// Create patterns of quoted regex
		for _, selectionPackage := range *packages {
			r, err := regexp.Compile("\\Q" + selectionPackage + "\\E")
			if err != nil {
				fmt.Printf("Error with quoted regex '%v' - %v", selectionPackage, err)
				panic(err)
			}
			regexPatterns = append(regexPatterns, r)
		}
	}





	availablePackages := testCatalog["packages"]

	for availablePackage, v:= range availablePackages.(map[string]interface{}) {

		for _, regexPackage := range regexPatterns {
			if regexPackage.MatchString(availablePackage) {
				availableClasses := v.([]interface{})				

				selectClasses(testCatalog, testSelection, availableClasses, stream)

				break
			}
		}	
	}
}

func selectClasses(testCatalog TestCatalog, testSelection *TestSelection, availableClasses []interface{}, stream string) {

	definedClasses := testCatalog["classes"].(map[string]interface{})

	if definedClasses == nil {
		return
	}


	for _, ac := range availableClasses {
		definedClass := definedClasses[ac.(string)].(map[string]interface{})

		if definedClass == nil {
			continue
		}

		appendClass(testSelection, definedClass, stream)

		bundle := definedClass["bundle"].(string)
		name   := definedClass["name"].(string)

		fmt.Printf("Selected test class '%v/%v'\n", bundle, name)
	}

}

func appendClass(testSelection *TestSelection, appendClass map[string]interface{}, stream string) {

	bundle := appendClass["bundle"].(string)
	name   := appendClass["name"].(string)	

	for _, selectedClass := range testSelection.Classes {
		if bundle == selectedClass.Bundle &&
		   name   == selectedClass.Class &&
		   stream == selectedClass.Stream {
			   return // already selected
		   }
	}	

	newSelectedClass := TestClass{
		Bundle: bundle,
		Class: name,
		Stream: stream,
	}

	testSelection.Classes = append(testSelection.Classes, newSelectedClass)

	fmt.Printf("Selected test class '%v/%v'\n", bundle, name)
}