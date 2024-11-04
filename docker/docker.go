package docker

import (
	"context"
	"fmt"
)

type RegistryAdapter interface {
	CheckAuth() error
	Pull(ctx context.Context, image, tag, targetPath string) error
	PushTar(context.Context, *Image) error
	Sync(context.Context, *Image) error
}

type Image struct {
	Name    string
	Tag     string
	TarPath string
}

func (i Image) String() string {
	return fmt.Sprintf("%s:%s, %s", i.Name, i.Tag, i.TarPath)
}
