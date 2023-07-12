/*
 * Copyright contributors to the Galasa project
 */
package utils

type MockUser struct {
	username string
}

func NewMockUser() *MockUser {
	user := new(MockUser)
	user.username = ""
	return user
}

func (env *MockUser) GetUsername() string {
	return env.username
}

func (env *MockUser) SetUsername(name string) {
	env.username = name
}
