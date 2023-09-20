/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */

package propertiesformatter

import (
	"fmt"
	"strings"
)

//Print in the following fashion:
// namespace	name	    value
// framework	property1	value1
// framework	property2	value2
// Total:1

// -----------------------------------------------------
// PropertyFormatter - implementations can take a collection of properties results
// and turn them into a string for display to the user.
const (
	HEADER_PROPERTY_NAMESPACE = "Namespace"
	HEADER_PROPERTY_NAME      = "Name"
	HEADER_PROPERTY_VALUE     = "Value"
)

type FormattableProperty struct {
	Namespace string
	Name      string
	Value     string
}

func NewFormattableProperty() FormattableProperty {
	this := FormattableProperty{}
	return this
}

type PropertyFormatter interface {
	FormatProperties(propertyResults []FormattableProperty) (string, error)
	GetName() string
}



// -----------------------------------------------------
// Functions for tables
func calculateMaxLengthOfEachColumn(table [][]string) []int {
	columnLengths := make([]int, len(table[0]))
	for _, row := range table {
		for i, val := range row {
			if len(val) > columnLengths[i] {
				columnLengths[i] = len(val)
			}
		}
	}
	return columnLengths
}

func writeFormattedTableToStringBuilder(table [][]string, buff *strings.Builder, columnLengths []int) {
	for _, row := range table {
		for column, val := range row {

			// For every column except the last one, add spacing.
			if column < len(row)-1 {
				// %-*s : variable space-padding length, padding is on the right.
				buff.WriteString(fmt.Sprintf("%-*s", columnLengths[column], val))
				buff.WriteString(" ")
			} else {
				buff.WriteString(val)
			}
		}
		buff.WriteString("\n")
	}
}