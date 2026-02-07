package avatar

import (
	"io"
)

type Store interface {
	Put(userID string, reader io.Reader) (avatar string, err error)
	Get(avatar string) (reader io.ReadCloser, size int64, err error)
	ID(avatar string) (id string)
	Remove(avatar string) error
}
