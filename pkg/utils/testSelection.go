/*
*  Copyright contributors to the Galasa project
 */

package utils

import (
	"log"
	"regexp"
	"strings"

	galasaErrors "github.com/galasa.dev/cli/pkg/errors"
	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/spf13/cobra"
)

type TestSelectionFlags struct {
	bundles     *[]string
	packages    *[]string
	tests       *[]string
	tags        *[]string
	classes     *[]string
	stream      string
	regexSelect *bool
}

type TestSelection struct {
	Classes []TestClass
}

type TestClass struct {
	Bundle string
	Class  string
	Stream string
}

func AddCommandFlags(command *cobra.Command, flags *TestSelectionFlags) {
	flags.packages = command.Flags().StringSlice("package", make([]string, 0), "packages of which tests will be selected from, packages are selected if the name contains this string, or if --regex is specified then matches the regex")
	flags.bundles = command.Flags().StringSlice("bundle", make([]string, 0), "bundles of which tests will be selected from, bundles are selected if the name contains this string, or if --regex is specified then matches the regex")
	flags.tests = command.Flags().StringSlice("test", make([]string, 0), "test names which will be selected if the name contains this string, or if --regex is specified then matches the regex")
	flags.tags = command.Flags().StringSlice("tag", make([]string, 0), "tags of which tests will be selected from, tags are selected if the name contains this string, or if --regex is specified then matches the regex")
	flags.classes = command.Flags().StringSlice("class", make([]string, 0), "test class names, for building a portfolio when a stream/test catalog is not available")
	command.Flags().StringVarP(&flags.stream, "stream", "s", "", "test stream to extract the tests from")
	flags.regexSelect = command.Flags().Bool("regex", false, "Test selection is performed by using regex")
}

func AreSelectionFlagsProvided(flags *TestSelectionFlags) bool {
	if len(*flags.bundles) > 0 {
		return true
	}

	if len(*flags.packages) > 0 {
		return true
	}

	if len(*flags.tests) > 0 {
		return true
	}

	if len(*flags.tags) > 0 {
		return true
	}

	if len(*flags.classes) > 0 {
		return true
	}

	if flags.stream != "" {
		return true
	}

	return false
}

func SelectTests(apiClient *galasaapi.APIClient, flags *TestSelectionFlags) TestSelection {

	var testCatalog TestCatalog
	if flags.stream != "" {
		availableStreams := FetchTestStreams(apiClient)

		err := ValidateStream(availableStreams, flags.stream)
		if err != nil {
			panic(err)
		}

		testCatalog, err = FetchTestCatalog(apiClient, flags.stream)
		if err != nil {
			panic(err)
		}
		log.Println("Test catalog retrieved")
	}

	testSelection := TestSelection{Classes: make([]TestClass, 0)}

	if flags.stream == "" {
		if len(*flags.packages) > 0 || len(*flags.bundles) > 0 || len(*flags.tests) > 0 {
			err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_STREAM_FLAG_REQUIRED)
			panic(err)
		}
	}

	selectTestsByBundle(testCatalog, &testSelection, flags)
	selectTestsByPackage(testCatalog, &testSelection, flags)
	selectTestsByTest(testCatalog, &testSelection, flags)
	selectTestsByTag(testCatalog, &testSelection, flags)
	selectTestsByClass(testCatalog, &testSelection, flags)

	return testSelection
}

func selectTestsByBundle(testCatalog TestCatalog, testSelection *TestSelection, flags *TestSelectionFlags) {

	if len(*flags.bundles) < 1 {
		return
	}

	if testCatalog["classes"] == nil {
		return
	}

	regexPatterns := convertRegex(flags.bundles, *flags.regexSelect)
	availableClasses := testCatalog["classes"].(map[string]interface{})

	for _, oclassDef := range availableClasses {
		classDef := oclassDef.(map[string]interface{})
		for _, regexBundle := range *regexPatterns {
			if regexBundle.MatchString(classDef["bundle"].(string)) {
				selectClassByCatalog(testSelection, classDef, flags)
				break
			}
		}
	}
}

func selectTestsByPackage(testCatalog TestCatalog, testSelection *TestSelection, flags *TestSelectionFlags) {

	if len(*flags.packages) < 1 {
		return
	}

	if testCatalog["classes"] == nil {
		return
	}

	regexPatterns := convertRegex(flags.packages, *flags.regexSelect)
	availableClasses := testCatalog["classes"].(map[string]interface{})

	for _, oclassDef := range availableClasses {
		classDef := oclassDef.(map[string]interface{})
		for _, regexBundle := range *regexPatterns {
			if regexBundle.MatchString(classDef["package"].(string)) {
				selectClassByCatalog(testSelection, classDef, flags)
				break
			}
		}
	}
}

func selectTestsByTest(testCatalog TestCatalog, testSelection *TestSelection, flags *TestSelectionFlags) {

	if len(*flags.tests) < 1 {
		return
	}

	if testCatalog["classes"] == nil {
		return
	}

	regexPatterns := convertRegex(flags.tests, *flags.regexSelect)
	availableClasses := testCatalog["classes"].(map[string]interface{})

	for _, oclassDef := range availableClasses {
		classDef := oclassDef.(map[string]interface{})
		for _, regexBundle := range *regexPatterns {
			if regexBundle.MatchString(classDef["name"].(string)) {
				selectClassByCatalog(testSelection, classDef, flags)
				break
			}
		}
	}
}

func selectTestsByClass(testCatalog TestCatalog, testSelection *TestSelection, flags *TestSelectionFlags) {

	if len(*flags.classes) < 1 {
		return
	}

	for _, class := range *flags.classes {
		pos := strings.Index(class, "/")
		if pos < 1 {
			err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CLASS_FORMAT, class)
			panic(err)
		}

		bundle := class[:pos]
		name := class[pos+1:]

		if name == "" {
			err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_CLASS_NAME_BLANK, class)
			panic(err)
		}

		selectClass(testSelection, bundle, name, flags)
	}
}

func selectTestsByTag(testCatalog TestCatalog, testSelection *TestSelection, flags *TestSelectionFlags) {

	if len(*flags.tags) < 1 {
		return
	}

	if testCatalog["classes"] == nil {
		return
	}

	regexPatterns := convertRegex(flags.tags, *flags.regexSelect)
	availableClasses := testCatalog["classes"].(map[string]interface{})

classSearch:
	for _, oclassDef := range availableClasses {
		classDef := oclassDef.(map[string]interface{})
		oclassTags := classDef["tags"]
		if oclassTags != nil {
			for _, itag := range oclassTags.([]interface{}) {
				tag := itag.(string)

				for _, regexBundle := range *regexPatterns {
					if regexBundle.MatchString(tag) {
						selectClassByCatalog(testSelection, classDef, flags)
						continue classSearch
					}
				}
			}
		}
	}
}

func selectClassByCatalog(testSelection *TestSelection, appendClass map[string]interface{}, flags *TestSelectionFlags) {

	bundle := appendClass["bundle"].(string)
	name := appendClass["name"].(string)

	selectClass(testSelection, bundle, name, flags)
}

func selectClass(testSelection *TestSelection, bundle string, name string, flags *TestSelectionFlags) {

	for _, selectedClass := range testSelection.Classes {
		if bundle == selectedClass.Bundle &&
			name == selectedClass.Class &&
			flags.stream == selectedClass.Stream {
			return // already selected
		}
	}

	newSelectedClass := TestClass{
		Bundle: bundle,
		Class:  name,
		Stream: flags.stream,
	}

	testSelection.Classes = append(testSelection.Classes, newSelectedClass)

	log.Printf("    Selected test class '%v/%v'\n", bundle, name)
}

func convertRegex(patterns *[]string, regexSelect bool) *[]*regexp.Regexp {

	regexPatterns := make([]*regexp.Regexp, 0)
	if regexSelect {
		// Create patterns of actual regex
		for _, selection := range *patterns {
			r, err := regexp.Compile(selection)
			if err != nil {
				err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SELECTION_REGEX_ERROR, selection, err.Error())
				panic(err)
			}
			regexPatterns = append(regexPatterns, r)
		}
	} else {
		// Create patterns of quoted regex
		for _, selection := range *patterns {
			r, err := regexp.Compile("\\Q" + selection + "\\E")
			if err != nil {
				err := galasaErrors.NewGalasaError(galasaErrors.GALASA_ERROR_SELECTION_REGEX_QUOTED_ERROR, selection, err.Error())
				panic(err)
			}
			regexPatterns = append(regexPatterns, r)
		}
	}

	return &regexPatterns
}
