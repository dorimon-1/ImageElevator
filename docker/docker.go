package docker

import "context"

type RegistryAdapter interface {
	CheckAuth() error
	Pull(ctx context.Context, image, tag, targetPath string) error
	PushTar(context.Context, *Image) error
}
