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

var dockerRunner docker.Docker

func init() {
	containerConfig := config.DockerConfig()
	dockerRunner = &docker.Container{
		Config: &containerConfig,
	}
}

func TestCheckAuth(t *testing.T) {
	if err := dockerRunner.CheckAuth(); err != nil {
		t.Errorf("failed: %v", err)
	}
}

func TestPull(t *testing.T) {
	tarPath := ".."

	if err := dockerRunner.Pull(imageName, tag, tarPath); err != nil {
		t.Fatalf("pulling image: %s", err)
	}
}

func TestPush(t *testing.T) {
	tarPath := fmt.Sprintf("../%s-%s", imageName, tag)

	if err := dockerRunner.PushTar(tarPath, imageName, tag); err != nil {
		t.Errorf("failed pushing: %v", err)
	}

	if err := os.Remove(tarPath); err != nil {
		t.Errorf("failed to delete file: %s", err)
	}
}
