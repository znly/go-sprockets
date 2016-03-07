package sprockets

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/znly/go-sprockets/types"
)

type requireTree struct {
	Path    string
	BaseDir string
}

type requireDirectory struct {
	Path    string
	BaseDir string
}

type requireFile struct {
	Path    string
	BaseDir string
}

// GetList is needed for RequireInterface
func (rt *requireTree) GetList(extInfo *types.ExtensionInfo) (requiredFiles []string, lastModified int64, err error) {
	finalPath := rt.Path
	if !filepath.IsAbs(finalPath) {
		finalPath = filepath.Join(rt.BaseDir, rt.Path)
		if newPath, err := filepath.EvalSymlinks(finalPath); err == nil {
			finalPath = newPath
		}
	}
	alterExts := extInfo.AlterExts
	if alterExts.Len() == 0 {
		return nil, 0, errors.New("Extension Format unknown for require:" + extInfo.CurrentExtension)
	}
	err = filepath.Walk(finalPath, func(walkpath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if modified := f.ModTime().Unix(); lastModified < modified {
			lastModified = modified
		}
		if f.IsDir() {
			return nil
		}
		curExt := filepath.Ext(walkpath)
		if curExt == extInfo.CurrentExtension {
			requiredFiles = append(requiredFiles, walkpath)
			return nil
		}
		for e := alterExts.Front(); e != nil; e = e.Next() {
			if e.Value == curExt {
				requiredFiles = append(requiredFiles, walkpath)
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	sort.Strings(requiredFiles)
	return
}

// GetList is needed for RequireInterface
func (rd *requireDirectory) GetList(extInfo *types.ExtensionInfo) (requiredFiles []string, lastModified int64, err error) {
	finalPath := rd.Path
	if !filepath.IsAbs(finalPath) {
		finalPath = filepath.Join(rd.BaseDir, rd.Path)
		if newPath, err := filepath.EvalSymlinks(finalPath); err == nil {
			finalPath = newPath
		}
	}
	f, err := os.Stat(finalPath)
	if err != nil {
		return nil, 0, err
	}
	if modified := f.ModTime().Unix(); lastModified < modified {
		lastModified = modified
	}
	alterExts := extInfo.AlterExts
	if alterExts.Len() == 0 {
		return nil, 0, errors.New("Extension Format unknown for require:" + extInfo.CurrentExtension)
	}
	for e := alterExts.Front(); e != nil; e = e.Next() {
		files, err := filepath.Glob(filepath.Join(finalPath, "*"+e.Value))
		if err != nil {
			return nil, 0, err
		}
		for _, file := range files {
			f, err := os.Stat(file)
			if err != nil {
				return nil, 0, err
			}
			if modified := f.ModTime().Unix(); lastModified < modified {
				lastModified = modified
			}
		}

		requiredFiles = append(requiredFiles, files...)
	}
	sort.Strings(requiredFiles)
	return

}

// GetList is needed for RequireInterface
func (rf *requireFile) GetList(extInfo *types.ExtensionInfo) ([]string, int64, error) {
	assetPath, _, err := resolvePath(extInfo, rf.Path, rf.BaseDir)
	if err != nil {
		return nil, 0, err
	}
	f, err := os.Stat(assetPath)
	if err != nil {
		return nil, 0, err
	}
	return []string{assetPath}, f.ModTime().Unix(), nil
}
