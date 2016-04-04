package sprockets

import (
	"path/filepath"

	"github.com/znly/go-sprockets/assetscache"
	"github.com/znly/go-sprockets/types"
)

// New creates a new Sprocket pipeline
// if publicPath is an empty string, Asset will be compiled and cached in memory
// if publicPath is not empty Files will be served from the public path and builded only if necessary
func New(assetsPath, publicPath string) (s *Sprocket, err error) {
	s = &Sprocket{}
	s.assetsPath, err = filepath.Abs(assetsPath)
	if err != nil {
		return nil, err
	}
	s.defaultExtInfo = newExtensionInfo("")
	s.PushFrontDefaultPath(s.assetsPath)
	s.extInfos = make(map[string]*types.ExtensionInfo)
	s.assetsCache = assetscache.New()
	if len(publicPath) == 0 {
		return
	}
	s.publicPath, err = filepath.Abs(publicPath)
	if err != nil {
		return nil, err
	}
	return
}

func (s *Sprocket) getAsset(assetPath string, forceRebuild bool) ([]byte, error) {
	realAssetPath, extInfo, err := s.resolvePath(assetPath, "", forceRebuild)
	if err != nil {
		return nil, err
	}
	var cacheKey *assetscache.AssetCacheKey
	if cacheKey, err = s.assetsCache.GenerateCacheKey(realAssetPath); err != nil {
		return nil, err
	}
	if forceRebuild == false {
		if cachedfullContent, err := s.assetsCache.GetFullCache(cacheKey); cachedfullContent != nil || err != nil {
			return cachedfullContent, err
		}
	}
	fullContent, content, requires, err := s.readAsset(realAssetPath, extInfo, forceRebuild)
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
	if forceRebuild == true {
		return fullContent, s.writeToPublic(assetPath, fullContent)
	}
	go func() {
		s.writeToPublic(assetPath, fullContent)
	}()
	return fullContent, nil
}

// GetAsset will return the asset full content (with all its requirement) or an error if an error occured
func (s *Sprocket) GetAsset(assetPath string) ([]byte, error) {
	return s.getAsset(assetPath, false)
}
