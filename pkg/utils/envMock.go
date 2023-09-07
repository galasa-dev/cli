/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package utils

type MockEnv struct {
	EnvVars map[string]string
}

func NewMockEnv() *MockEnv {
	env := new(MockEnv)
	env.EnvVars = make(map[string]string, 0)
	return env
}

func (env *MockEnv) GetEnv(propertyName string) string {
	return env.EnvVars[propertyName]
}

func (env *MockEnv) SetEnv(propertyName string, value string) {
	env.EnvVars[propertyName] = value
}
