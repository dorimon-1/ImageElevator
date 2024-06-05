package containers

import (
	"context"
	"fmt"

	"github.com/Kjone1/imageElevator/config"
	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/signature"
	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

// CheckAuth is used to check authentication only to DockerHub
func CheckAuth(config *config.ContainerConfig) error {
	return docker.CheckAuth(
		context.Background(),
		config.SystemContext,
		config.SystemContext.DockerAuthConfig.Username,
		config.SystemContext.DockerAuthConfig.Password,
		config.Registry,
	)
}

func Login(repository, tag string) (types.ImageReference, error) {
	return parseDocker(repository, tag)
}

func PushTar(tarPath, tag string, config *config.ContainerConfig) error {
	dstRef, err := parseDocker(config.RepoURL, tag)
	if err != nil {
		return err
	}

	srcRef, err := parseTar(tarPath)
	if err != nil {
		return err
	}

	policyCtx, err := signature.NewPolicyContext(&signature.Policy{
		Default: []signature.PolicyRequirement{
			signature.NewPRInsecureAcceptAnything(),
		},
	})
	if err != nil {
		return err
	}

	defer func() { _ = policyCtx.Destroy() }()

	if _, err := copy.Image(context.Background(), policyCtx, dstRef, srcRef, &copy.Options{
		DestinationCtx: config.SystemContext,
		SourceCtx:      config.SystemContext,
	}); err != nil {
		return fmt.Errorf("coping source image to destination repository: %s", err)
	}

	return nil
}

func parseTar(path string) (types.ImageReference, error) {
	ref, err := alltransports.ParseImageName(fmt.Sprintf("docker-archive:%s", path))
	if err != nil {
		return nil, fmt.Errorf("parsing %s to image name: %s", path, err)
	}
	return ref, nil

}

func parseDocker(repository, tag string) (types.ImageReference, error) {
	ref, err := alltransports.ParseImageName(fmt.Sprintf("docker://%s:%s", repository, tag))
	if err != nil {
		return nil, fmt.Errorf("parsing repository on login: %s", err)
	}
	return ref, nil
}
