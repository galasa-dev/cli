/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package properties

import (
	"sort"
	"strings"

	"github.com/galasa.dev/cli/pkg/galasaapi"
	"github.com/galasa.dev/cli/pkg/propertiesformatter"
)

func orderFormattableProperties(formattableProperty []propertiesformatter.FormattableProperty) []propertiesformatter.FormattableProperty {

	//order by property name
	sort.SliceStable(formattableProperty, func(i, j int) bool {
		return formattableProperty[i].Name < formattableProperty[j].Name
	})

	return formattableProperty
}

func FormattablePropertyFromGalasaApi(properties []galasaapi.CpsProperty) []propertiesformatter.FormattableProperty {
	var formattableProperty []propertiesformatter.FormattableProperty

	for _, property := range properties {
		//Get the data for each property
		newFormattableProperty := getCpsPropertyData(property)
		formattableProperty = append(formattableProperty, newFormattableProperty)
	}

	orderedFormattableProperties := orderFormattableProperties(formattableProperty)

	return orderedFormattableProperties
}

func getCpsPropertyData(property galasaapi.CpsProperty) propertiesformatter.FormattableProperty {
	newFormattableProperty := propertiesformatter.NewFormattableProperty()

	//GetName() returns the full name which is in the form namespace.name
	fullName := property.GetName()
	firstDotIndex := strings.Index(fullName, ".")

	//for namsepace capture everything before the first dot
	newFormattableProperty.Namespace = fullName[:firstDotIndex]
	//for name, capture everything after the first dot
	newFormattableProperty.Name = fullName[firstDotIndex+1:]
	newFormattableProperty.Value = property.GetValue()

	return newFormattableProperty
}