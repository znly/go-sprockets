package assetscache

import (
	"errors"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/znly/go-sprockets/dependencygraph"
	"github.com/znly/go-sprockets/types"
)

var (
	errMustRebuildCache = errors.New("")
)

// AssetsCache structure use for caching assets
type AssetsCache struct {
	mutex *sync.RWMutex
	cache map[string]*assetLru
}

// AssetCacheKey structure use as key for caching assets
type AssetCacheKey struct {
	AssetPath string
	Key       int64
}

type assetCache struct {
	Requires    []types.RequireInterface
	FullContent []byte
	Content     []byte
	ExtInfo     *types.ExtensionInfo
	LastWrite   int64
}

// New return a new AssetCache structure
func New() (a *AssetsCache) {
	a = &AssetsCache{
		mutex: &sync.RWMutex{},
		cache: make(map[string]*assetLru),
	}
	return
}

// GenerateCacheKey will generate a new cache key base on the asset time stamp
// return an error if os.Stat failed.
func (a *AssetsCache) GenerateCacheKey(assetPath string) (*AssetCacheKey, error) {
	info, err := os.Stat(assetPath)
	if err != nil {
		return nil, err
	}
	key := info.ModTime().Unix()
	return &AssetCacheKey{assetPath, key}, nil
}

func (a *AssetsCache) readFromCache(key *AssetCacheKey) *assetCache {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	if assetCaches, hit := a.cache[key.AssetPath]; hit {
		if cache, hit := assetCaches.Get(key.Key); hit {
			return cache
		}
	}
	return nil
}

// ReadFromCache will return the current cache content
// if retHit is false, the cache is empty or outdated for this asset
func (a *AssetsCache) ReadFromCache(key *AssetCacheKey) (content []byte, requires []types.RequireInterface, fullContent []byte, extInfo *types.ExtensionInfo, retHit bool) {
	cache := a.readFromCache(key)
	if cache == nil {
		return
	}
	retHit = true
	requires = cache.Requires
	content = cache.Content
	fullContent = cache.FullContent
	extInfo = cache.ExtInfo
	return
}

// WriteToCache will write the content of an asset into the cache
func (a *AssetsCache) WriteToCache(key *AssetCacheKey, fullContent, content []byte, requires []types.RequireInterface, ExtInfo *types.ExtensionInfo) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	var assetCaches *assetLru
	var hit bool
	if assetCaches, hit = a.cache[key.AssetPath]; !hit {
		assetCaches = newAssetLru(5)
		a.cache[key.AssetPath] = assetCaches
	}
	assetCaches.Add(key.Key, &assetCache{requires, fullContent, content, ExtInfo, time.Now().Unix()})
}

// GetFullCache will return the full content of a Cache if it s available and not outdated
func (a *AssetsCache) GetFullCache(key *AssetCacheKey) ([]byte, error) {
	cache := a.readFromCache(key)
	if cache == nil || cache.FullContent == nil {
		return nil, nil
	}
	graph := dependencygraph.Graph{}
	_, err := graph.Walk(key.AssetPath, func(curPath, parentPath string, g *dependencygraph.Graph) error {
		curKey, err := a.GenerateCacheKey(curPath)
		if err != nil {
			return err
		}
		if curKey.Key > cache.LastWrite {
			return errMustRebuildCache
		}
		curCache := a.readFromCache(curKey)
		if curCache == nil {
			return errMustRebuildCache
		}
		for _, r := range curCache.Requires {
			requiredFiles, lastModified, err := r.GetList(curCache.ExtInfo)
			if err != nil {
				return err
			}
			if lastModified > cache.LastWrite {
				return errMustRebuildCache
			}
			selfIndex := sort.SearchStrings(requiredFiles, curPath)
			if selfIndex < len(requiredFiles) && requiredFiles[selfIndex] == curPath {
				requiredFiles = append(requiredFiles[:selfIndex], requiredFiles[selfIndex+1:]...)
			}
			g.AddChildrens(curPath, requiredFiles...)
		}
		return nil
	})
	if err == errMustRebuildCache {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return cache.FullContent, nil
}
