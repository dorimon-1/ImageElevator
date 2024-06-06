package containers_test

import (
	"testing"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/containers"
)

const (
	tag       = "3"
	imageName = "dor4420"
)

func TestPush(t *testing.T) {
	config.LoadConfig()
	containerConfig := config.Config.ContainerConfig
	tarPath := "../alpine.tar"

	if err := containers.PushTar(tarPath, imageName, tag, containerConfig); err != nil {
		t.Errorf("failed pushing: %v", err)
	}

}

func TestLogin(t *testing.T) {
	config.LoadConfig()
	containerConfig := config.Config.ContainerConfig
	_, err := containers.Login(containerConfig.Registry, containerConfig.Repository, imageName, tag)
	if err != nil {
		t.Errorf("failed: %v", err)
	}

}
