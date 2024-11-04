package ftp

type FTPClient interface {
	Pull(files ...string) ([]string, error)
	List(path string, pattern string, bannedFiles map[string]bool) ([]string, error)
	Close() error
}
