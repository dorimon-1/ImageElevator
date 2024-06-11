package docker

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

type Container struct {
	*config.RegistryConfiguration
}

func (c *Container) CheckAuth() error {
	return docker.CheckAuth(
		context.Background(),
		c.SystemContext,
		c.SystemContext.DockerAuthConfig.Username,
		c.SystemContext.DockerAuthConfig.Password,
		c.Registry,
	)
}

func (c *Container) Pull(image, tag, targetPath string) error {
	imgRef, err := parseDocker(c.Registry, c.Repository, image, tag)
	if err != nil {
		return err
	}

	dstRef, err := parseTar(fmt.Sprintf("%s/%s-%s", targetPath, image, tag))
	if err != nil {
		return err
	}

	if err := copyImage(imgRef, dstRef, c.SystemContext); err != nil {
		return err
	}

	return nil
}

func (c *Container) PushTar(tarPath, imageName, tag string) error {
	dstRef, err := parseDocker(c.Registry, c.Repository, imageName, tag)
	if err != nil {
		return err
	}

	srcRef, err := parseTar(tarPath)
	if err != nil {
		return err
	}

	if err := copyImage(srcRef, dstRef, c.SystemContext); err != nil {
		return err
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

func parseDocker(registry, repository, imageName, tag string) (types.ImageReference, error) {
	ref, err := alltransports.ParseImageName(fmt.Sprintf("docker://%s/%s/%s:%s", registry, repository, imageName, tag))
	if err != nil {
		return nil, fmt.Errorf("parsing repository on login: %s", err)
	}

	return ref, nil
}

func copyImage(srcRef, dstRef types.ImageReference, sysCtx *types.SystemContext) error {
	policyCtx, err := signature.NewPolicyContext(&signature.Policy{
		Default: []signature.PolicyRequirement{
			signature.NewPRInsecureAcceptAnything(),
		},
	})
	if err != nil {
		return err
	}

	defer func() { _ = policyCtx.Destroy() }()

	_, err = copy.Image(context.Background(), policyCtx, dstRef, srcRef, &copy.Options{
		SourceCtx:      sysCtx,
		DestinationCtx: sysCtx,
	})

	return err
}
