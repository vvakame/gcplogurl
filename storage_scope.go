package gcplogurl

import (
	"io"
	"net/url"
)

// StorageScope means scope about logs.
type StorageScope interface {
	isStorageScope()
	marshalURL(w io.Writer)
}

var _ StorageScope = StorageScopeProject
var _ StorageScope = (*StorageScopeStorage)(nil)

// StorageScopeProject use scope by project.
const StorageScopeProject = storageScopeProject(0)

type storageScopeProject int

func (storageScopeProject) isStorageScope() {}

func (storageScopeProject) marshalURL(w io.Writer) {
	_, _ = w.Write([]byte(";storageScope=project"))
}

// StorageScopeStorage use scope by storage.
type StorageScopeStorage struct {
	Src []string
}

func (s *StorageScopeStorage) isStorageScope() {}

func (s *StorageScopeStorage) marshalURL(w io.Writer) {
	_, _ = w.Write([]byte(";storageScope=storage"))
	for _, src := range s.Src {
		_, _ = w.Write([]byte(","))
		_, _ = w.Write([]byte(url.QueryEscape(src)))
	}
}
