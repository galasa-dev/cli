/*
 * Copyright contributors to the Galasa project
 */
package utils

type MockEnv struct {
	EnvVars  map[string]string
	username string
}

func NewMockEnv() *MockEnv {
	env := new(MockEnv)
	env.EnvVars = make(map[string]string, 0)
	env.username = ""
	return env
}

func (env *MockEnv) GetEnv(propertyName string) string {
	return env.EnvVars[propertyName]
}

func (env *MockEnv) SetEnv(propertyName string, value string) {
	env.EnvVars[propertyName] = value
}

func (env *MockEnv) GetUsername() string {
	return env.username
}

func (env *MockEnv) SetUsername(name string) {
	env.username = name
}
