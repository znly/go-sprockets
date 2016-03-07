package sprockets

import (
	"path/filepath"

	"github.com/znly/go-sprockets/assetscache"
	"github.com/znly/go-sprockets/types"
)

// New creates a new Sprocket pipeline
func New(assetDir string) (s *Sprocket, err error) {
	s = &Sprocket{}
	absAssetDir, err := filepath.Abs(assetDir)
	if err != nil {
		return nil, err
	}
	s.assetPath = absAssetDir
	s.defaultExtInfo = newExtensionInfo("")
	s.PushFrontDefaultPath(s.assetPath)
	s.extInfos = make(map[string]*types.ExtensionInfo)
	s.assetsCache = assetscache.New()
	return
}

// GetAsset will return the asset full content (with all its requirement) or an error if an error occured
func (s *Sprocket) GetAsset(assetPath string) ([]byte, error) {
	realAssetPath, extInfo, err := s.resolvePath(assetPath, "")
	if err != nil {
		return nil, err
	}
	var cacheKey *assetscache.AssetCacheKey
	if cacheKey, err = s.assetsCache.GenerateCacheKey(realAssetPath); err != nil {
		return nil, err
	}

	if cachedfullContent, err := s.assetsCache.GetFullCache(cacheKey); cachedfullContent != nil || err != nil {
		return cachedfullContent, err
	}
	fullContent, content, requires, err := s.readAsset(realAssetPath, extInfo)
	if err != nil {
		return nil, err
	}
	if extInfo.BundleCompiler != nil {
		fullContent, err = extInfo.BundleCompiler.Process(fullContent, realAssetPath)
		if err != nil {
			return nil, err
		}
	}
	for _, f := range extInfo.PostCompileContentTreatment {
		fullContent, err = f.Process(fullContent, realAssetPath)
		if err != nil {
			return nil, err
		}
	}
	s.assetsCache.WriteToCache(cacheKey, fullContent, content, requires, extInfo)
	return fullContent, nil
}
