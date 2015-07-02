package core

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"time"
)

func GetFileLocation(node *Node) string {
	strUuid := node.Uuid.CleanString()

	return fmt.Sprintf("%s/%s/%s.bin", strUuid[0:2], strUuid[2:4], strUuid[4:])
}

func NewSecureFs(fs afero.Fs, path string) *SecureFs {
	return &SecureFs{
		fs:   fs,
		path: path,
	}
}

type SecureFs struct {
	fs   afero.Fs
	path string
}

func (s *SecureFs) securePath(name string) string {
	return fmt.Sprintf("%s/%s", s.path, name)
}

func (s *SecureFs) Create(name string) (afero.File, error) {
	return s.fs.Create(s.securePath(name))
}

func (s *SecureFs) Mkdir(name string, perm os.FileMode) error {
	return s.fs.Mkdir(s.securePath(name), perm)
}

func (s *SecureFs) MkdirAll(path string, perm os.FileMode) error {
	return s.fs.MkdirAll(s.securePath(path), perm)
}

func (s *SecureFs) Open(name string) (afero.File, error) {
	return s.fs.Open(s.securePath(name))
}

func (s *SecureFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	return s.fs.OpenFile(s.securePath(name), flag, perm)
}

func (s *SecureFs) Remove(name string) error {
	return s.fs.Remove(s.securePath(name))
}

func (s *SecureFs) RemoveAll(path string) error {
	return s.fs.RemoveAll(path)
}

func (s *SecureFs) Rename(oldname, newname string) error {
	return s.fs.Rename(s.securePath(oldname), s.securePath(newname))
}

func (s *SecureFs) Stat(name string) (os.FileInfo, error) {
	return s.fs.Stat(s.securePath(name))
}

func (s *SecureFs) Name() string {
	return "SecureFs"
}

func (s *SecureFs) Chmod(name string, mode os.FileMode) error {
	return s.fs.Chmod(s.securePath(name), mode)
}

func (s *SecureFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return s.fs.Chtimes(s.securePath(name), atime, mtime)
}
