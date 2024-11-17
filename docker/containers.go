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
	"github.com/rs/zerolog/log"
)

type Container struct {
	*config.RegistryConfiguration
	SyncRegistryConfigruation []config.RegistryConfiguration
}

func NewRegistry(config *config.RegistryConfiguration, syncConfig []config.RegistryConfiguration) RegistryAdapter {
	return &Container{
		RegistryConfiguration:     config,
		SyncRegistryConfigruation: syncConfig,
	}
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

func (c *Container) Pull(ctx context.Context, image, tag, targetPath string) error {
	imgRef, err := parseDocker(c.Registry, c.Repository, image, tag)
	if err != nil {
		return err
	}

	dstRef, err := parseTar(fmt.Sprintf("%s/%s-%s.tar", targetPath, image, tag))
	if err != nil {
		return err
	}

	if err := copyImage(ctx, imgRef, dstRef, c.SystemContext, c.SystemContext); err != nil {
		return err
	}

	return nil
}

func (c *Container) PushTar(ctx context.Context, image *Image) error {
	dstRef, err := parseDocker(c.Registry, c.Repository, image.Name, image.Tag)
	if err != nil {
		log.Error().Msgf("error parsing destination registry %s, repository %s, image %s:%s", c.Registry, c.Repository, image.Name, image.Tag)
		return err
	}

	log.Info().Msgf("parsing image tar path: %s", image.TarPath)
	srcRef, err := parseTar(image.TarPath)
	if err != nil {
		return err
	}
	log.Info().Msgf("Pushing %s/%s:%s at %s to %s", c.Repository, image.Name, image.Tag, image.TarPath, c.Registry)
	if err := copyImage(ctx, srcRef, dstRef, c.SystemContext, c.SystemContext); err != nil {
		log.Error().Msgf("error pushing image destination registry %s, repository %s, image %s:%s", c.Registry, c.Repository, image.Name, image.Tag)
		return err
	}

	return nil
}

func (c *Container) Sync(ctx context.Context, image *Image) error {
	srcRef, err := parseDocker(c.Registry, c.Repository, image.Name, image.Tag)
	if err != nil {
		return err
	}

	for _, rc := range c.SyncRegistryConfigruation {
		dstRef, err := parseDocker(rc.Repository, rc.Repository, image.Name, image.Tag)
		if err != nil {
			log.Error().Msgf("error parsing docker registry image %s:%s, registry %s: %s", image.Name, image.Tag, rc.Registry, err)
			continue
		}

		if err := copyImage(ctx, srcRef, dstRef, c.SystemContext, rc.SystemContext); err != nil {
			log.Error().Msgf("error syncing image %s:%s to registry %s: %s", image.Name, image.Tag, rc.Registry, err)
		}
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

func copyImage(ctx context.Context, srcRef, dstRef types.ImageReference, srcSysCtx *types.SystemContext, dstSysCtx *types.SystemContext) error {
	policyCtx, err := signature.NewPolicyContext(&signature.Policy{
		Default: []signature.PolicyRequirement{
			signature.NewPRInsecureAcceptAnything(),
		},
	})
	if err != nil {
		return err
	}

	defer func() { _ = policyCtx.Destroy() }()

	_, err = copy.Image(ctx, policyCtx, dstRef, srcRef, &copy.Options{
		SourceCtx:      srcSysCtx,
		DestinationCtx: dstSysCtx,
	})

	return err
}
