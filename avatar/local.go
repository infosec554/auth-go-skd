package avatar

import (
	"io"
	"os"
	"path"
)

type LocalFS struct {
	Location string
}

func NewLocalFS(location string) *LocalFS {
	os.MkdirAll(location, 0o700)
	return &LocalFS{Location: location}
}

func (l *LocalFS) Put(userID string, reader io.Reader) (avatar string, err error) {

	avatar = userID + ".image"
	file, err := os.Create(path.Join(l.Location, avatar))
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = io.Copy(file, reader)
	return avatar, err
}

func (l *LocalFS) Get(avatar string) (reader io.ReadCloser, size int64, err error) {
	file, err := os.Open(path.Join(l.Location, avatar))
	if err != nil {
		return nil, 0, err
	}
	info, _ := file.Stat()
	return file, info.Size(), nil
}

func (l *LocalFS) ID(avatar string) (id string) {
	return avatar
}

func (l *LocalFS) Remove(avatar string) error {
	return os.Remove(path.Join(l.Location, avatar))
}
