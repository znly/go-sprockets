package sprockets

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/znly/go-sprockets/types"
)

func resolveExt(ei *types.ExtensionInfo, argAssetPath, argExt string) (string, string, bool) {
	if isFileExist(argAssetPath) {
		return argAssetPath, argExt, true
	}
	if ei.AlterExts.Find(argExt) != nil {
		for e := ei.AlterExts.Front(); e != nil; e = e.Next() {
			alterExt := e.Value
			if alterExt == argExt {
				continue
			}
			alterPath := strings.Replace(argAssetPath, argExt, alterExt, 1)
			if isFileExist(alterPath) {
				return alterPath, alterExt, true
			}
		}
	}
	for e := ei.AlterExts.Front(); e != nil; e = e.Next() {
		alterExt := e.Value
		if alterExt == argExt {
			continue
		}
		alterPath := argAssetPath + alterExt
		if isFileExist(alterPath) {
			return alterPath, alterExt, true
		}
	}
	return "", "", false
}

func resolvePath(ei *types.ExtensionInfo, assetPath string, baseDir string) (string, string, error) {
	ext := filepath.Ext(assetPath)
	//TODO: add a security to stay in certain directories (CHROOT LIKE)
	if strings.HasPrefix(assetPath, ".") {
		if baseDir == "" {
			return "", "", errCantFindAsset(assetPath)
		}
		assetPath = filepath.Join(baseDir, assetPath)
	}
	if strings.HasPrefix(assetPath, "/") {
		curAssetPath, curExt, ok := resolveExt(ei, assetPath, ext)
		if !ok {
			return "", "", errCantFindAsset(assetPath)
		}
		assetPath = curAssetPath
		ext = curExt
	} else {
		found := false
		for e := ei.Paths.Front(); e != nil; e = e.Next() {
			curAssetPath := filepath.Join(e.Value, assetPath)
			if curAssetPath, curExt, ok := resolveExt(ei, curAssetPath, ext); ok {
				found = true
				assetPath = curAssetPath
				ext = curExt
				break
			}
		}
		if !found && len(baseDir) > 0 {
			curAssetPath := filepath.Join(baseDir, assetPath)
			if curAssetPath, curExt, ok := resolveExt(ei, curAssetPath, ext); ok {
				found = true
				assetPath = curAssetPath
				ext = curExt
			}
		}
		if !found {
			return "", "", errCantFindAsset(assetPath)
		}
	}
	if newAssetPath, err := filepath.EvalSymlinks(assetPath); err == nil {
		assetPath = newAssetPath
	}
	return assetPath, ext, nil
}

// resolvePath
// Return Pathfound, Extension found, or an error
// It will search base on path then extension
func (s *Sprocket) resolvePath(assetPath string, baseDir string) (string, *types.ExtensionInfo, error) {
	var err error
	ext := filepath.Ext(assetPath)
	extInfo := s.getExtensionInfoOrDefault(ext)
	assetPath, ext, err = resolvePath(extInfo, assetPath, baseDir)
	if err != nil {
		return "", nil, err
	}
	return assetPath, s.getExtensionInfoOrDefault(ext), nil
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
