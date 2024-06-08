package docker

type Docker interface {
	CheckAuth() error
	Pull(image, tag, targetPath string) error
	PushTar(tarPath, imageName, tag string) error
}
