package containers_test

import (
	"testing"

	"github.com/Kjone1/imageElevator/config"
	"github.com/Kjone1/imageElevator/containers"
	"github.com/containers/image/v5/types"
)

const (
	username  = "dorsahar@icloud.com"
	password  = "dor4420!@"
	repo      = "docker.io/dor4420"
	imageName = "dor4420"
	tag       = "1"
	registry  = "docker.io"
)

var containerConfig *config.ContainerConfig = &config.ContainerConfig{
	Registry: registry,
	RepoURL:  repo,
	SystemContext: &types.SystemContext{
		DockerAuthConfig: &types.DockerAuthConfig{
			Username: username,
			Password: password,
		},
	},
}

func TestPush(t *testing.T) {
	tarPath := "../alpine.tar"

	if err := containers.PushTar(tarPath, imageName, tag, containerConfig); err != nil {
		t.Errorf("failed pushing: %v", err)
	}

}

func TestLogin(t *testing.T) {
	_, err := containers.Login(repo, imageName, tag)
	if err != nil {
		t.Errorf("failed: %v", err)
	}

}
