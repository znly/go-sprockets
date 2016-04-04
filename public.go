package sprockets

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	ErrNoPublicPathSet = errors.New("No public path set")
)

func (s *Sprocket) writeToPublic(assetPath string, FullContent []byte) error {
	if len(s.publicPath) == 0 {
		return ErrNoPublicPathSet
	}
	dirname := filepath.Dir(assetPath)
	fullDirPath := filepath.Join(s.publicPath, dirname)
	fullPath := filepath.Join(s.publicPath, assetPath)
	if err := os.MkdirAll(fullDirPath, os.FileMode(0750)); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fullPath, FullContent, os.FileMode(0640)); err != nil {
		return err
	}
	return nil
}

func (s *Sprocket) Generate(assetPath string) error {
	if len(s.publicPath) == 0 {
		return ErrNoPublicPathSet
	}
	_, err := s.getAsset(assetPath, true)
	return err
}
