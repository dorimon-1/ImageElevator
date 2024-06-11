package docker

type RegistryAdapter interface {
	CheckAuth() error
	Pull(image, tag, targetPath string) error
	PushTar(tarPath, imageName, tag string) error
}
