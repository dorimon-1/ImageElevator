package docker_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/docker"
)

const (
	tag       = "3"
	imageName = "dor4420"
)

func TestCheckAuth(t *testing.T) {
	containerConfig := config.ContainersConfig()
	if err := docker.CheckAuth(&containerConfig); err != nil {
		t.Errorf("failed: %v", err)
	}
}

func TestPull(t *testing.T) {
	containerConfig := config.ContainersConfig()
	tarPath := ".."

	if err := docker.Pull(&containerConfig, imageName, tag, tarPath); err != nil {
		t.Fatalf("pulling image: %s", err)
	}
}

func TestPush(t *testing.T) {
	containerConfig := config.ContainersConfig()
	tarPath := fmt.Sprintf("../%s-%s", imageName, tag)

	if err := docker.PushTar(tarPath, imageName, tag, &containerConfig); err != nil {
		t.Errorf("failed pushing: %v", err)
	}

	if err := os.Remove(tarPath); err != nil {
		t.Errorf("failed to delete file: %s", err)
	}
}
