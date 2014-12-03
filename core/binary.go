package gonode

import (
	"github.com/spf13/afero"
	"os"
	"time"
	"fmt"
)


func NewSecureFs(fs afero.Fs, path string) *secureFs {
	return &secureFs{
		fs: fs,
		path: path,
	}
}

type secureFs struct {
	fs   afero.Fs
	path string
}

func (s *secureFs) securePath(name string) string {
	return fmt.Sprintf("%s/%s", s.path, name)
}

func (s *secureFs) Create(name string) (afero.File, error) {
	return s.fs.Create(s.securePath(name))
}

func (s *secureFs) Mkdir(name string, perm os.FileMode) error {
	return s.fs.Mkdir(s.securePath(name), perm)
}

func (s *secureFs) MkdirAll(path string, perm os.FileMode) error {
	return s.fs.MkdirAll(s.securePath(path), perm)
}

func (s *secureFs) Open(name string) (afero.File, error) {
	return s.fs.Open(s.securePath(name))
}

func (s *secureFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return s.fs.OpenFile(s.securePath(name), flag, perm)
}

func (s *secureFs) Remove(name string) error {
	return s.fs.Remove(s.securePath(name))
}

func (s *secureFs) RemoveAll(path string) error {
	return s.fs.RemoveAll(path)
}

func (s *secureFs) Rename(oldname, newname string) error {
	return s.fs.Rename(s.securePath(oldname), s.securePath(newname))
}

func (s *secureFs) Stat(name string) (os.FileInfo, error) {
	return s.fs.Stat(s.securePath(name))
}

func (s *secureFs) Name() string {
	return "SecureFs"
}

func (s *secureFs) Chmod(name string, mode os.FileMode) error {
	return s.fs.Chmod(s.securePath(name), mode)
}

func (s *secureFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return s.fs.Chtimes(s.securePath(name), atime, mtime)
}
