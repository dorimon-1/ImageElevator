package mocks

import (
	"context"

	"github.com/Kjone1/imageElevator/docker"
	"github.com/stretchr/testify/mock"
)

type MockRegistry struct {
	mock.Mock
}

func (m *MockRegistry) CheckAuth() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRegistry) Pull(ctx context.Context, image string, tag string, targetPath string) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRegistry) PushTar(_ context.Context, _ *docker.Image) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRegistry) Sync(ctx context.Context, image *docker.Image) error {
	args := m.Called()
	return args.Error(0)
}
