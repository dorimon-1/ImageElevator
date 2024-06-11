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

var imageRegistry docker.RegistryAdapter

func init() {
	containerConfig := config.RegistryConfig()
	imageRegistry = &docker.Container{
		RegistryConfiguration: &containerConfig,
	}
}

func TestCheckAuth(t *testing.T) {
	if err := imageRegistry.CheckAuth(); err != nil {
		t.Errorf("failed: %v", err)
	}
}

func TestPull(t *testing.T) {
	tarPath := ".."

	if err := imageRegistry.Pull(imageName, tag, tarPath); err != nil {
		t.Fatalf("pulling image: %s", err)
	}
}

func TestPush(t *testing.T) {
	tarPath := fmt.Sprintf("../%s-%s", imageName, tag)

	if err := imageRegistry.PushTar(tarPath, imageName, tag); err != nil {
		t.Errorf("failed pushing: %v", err)
	}

	if err := os.Remove(tarPath); err != nil {
		t.Errorf("failed to delete file: %s", err)
	}
}
