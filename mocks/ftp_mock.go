package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockFTPClient struct {
	mock.Mock
}

func (m *MockFTPClient) Pull(files ...string) ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockFTPClient) List(path, pattern string) ([]string, error) {
	args := m.Called()
	strings, ok := args.Get(0).([]string)
	if !ok {
		return nil, args.Error(1)
	}
	return strings, args.Error(1)
}

func (m *MockFTPClient) Close() error {

	args := m.Called()
	return args.Error(0)
}
