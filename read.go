package sprockets

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/znly/go-sprockets/assetscache"
	"github.com/znly/go-sprockets/dependencygraph"
	"github.com/znly/go-sprockets/types"
)

func (s *Sprocket) readAsset(assetPath string, extInfo *types.ExtensionInfo, forceRebuild bool) (fullContent, content []byte, requires []types.RequireInterface, err error) {
	if extInfo.RequirePattern == nil {
		content, err = s.readAssetContent(assetPath, extInfo)
		return content, content, nil, err
	}
	graph := dependencygraph.Graph{}
	curAssetCache := make(map[string][]byte)
	dependencyList, err := graph.Walk(assetPath, func(curPath, parentPath string, g *dependencygraph.Graph) error {
		curRequires, curContent, curErr := s.readAssetWithDependencies(curPath, parentPath, forceRebuild)
		if curErr != nil {
			return curErr
		}
		if curPath == assetPath {
			content = curContent
			requires = curRequires
		}
		for _, r := range curRequires {
			requiredFiles, _, err := r.GetList(extInfo)
			if err != nil {
				return err
			}
			selfIndex := sort.SearchStrings(requiredFiles, curPath)
			if selfIndex < len(requiredFiles) && requiredFiles[selfIndex] == curPath {
				requiredFiles = append(requiredFiles[:selfIndex], requiredFiles[selfIndex+1:]...)
			}
			g.AddChildrens(curPath, requiredFiles...)
		}
		curAssetCache[curPath] = curContent
		return nil
	})
	if err != nil {
		return
	}
	for _, val := range dependencyList {
		fullContent = append(fullContent, byte('\n'))
		fullContent = append(fullContent, curAssetCache[val]...)
	}
	return
}

func (s *Sprocket) readAssetWithDependencies(argAssetPath, parentPath string, forceRebuild bool) (requires []types.RequireInterface, content []byte, err error) {
	var cacheKey *assetscache.AssetCacheKey
	assetPath, extInfo, err := s.resolvePath(argAssetPath, filepath.Dir(parentPath), forceRebuild)
	if err != nil {
		return nil, nil, err
	}
	if cacheKey, err = s.assetsCache.GenerateCacheKey(assetPath); err != nil {
		return nil, nil, errCantFindAsset(assetPath)
	}
	if forceRebuild == false {
		var hit bool
		content, requires, _, _, hit = s.assetsCache.ReadFromCache(cacheKey)
		if hit {
			return
		}
	}
	content, err = s.readAssetContent(assetPath, extInfo)
	if err != nil {
		return
	}
	if extInfo.RequirePattern == nil {
		return
	}
	header := extInfo.RequirePattern.Head.Find(content)
	if len(header) != 0 {
		newheader := header
		dirPath := filepath.Dir(assetPath)
		for _, f := range extInfo.HeaderTreatment {
			newheader, err = f.Process(newheader, assetPath)
			if err != nil {
				return
			}
		}
		for _, line := range bytes.Split(newheader, []byte("\n")) {
			ret := extInfo.RequirePattern.Require.FindAllSubmatch(line, -1)
			if ret == nil {
				continue
			}
			newheader = bytes.Replace(newheader, line, bytes.Replace(line, ret[0][1], []byte{}, 1), 1) //Remove the require from the newheader
			if bytes.Equal(ret[0][2], []byte("_tree")) {
				requires = append(requires, &requireTree{string(ret[0][3]), dirPath})
			} else if bytes.Equal(ret[0][2], []byte("_directory")) {
				requires = append(requires, &requireDirectory{string(ret[0][3]), dirPath})
			} else {
				requires = append(requires, &requireFile{string(ret[0][3]), dirPath})
			}
		}

		content = bytes.Replace(content, header, newheader, 1)
	}
	if extInfo.FileCompiler != nil {
		content, err = extInfo.FileCompiler.Process(content, assetPath)
		if err != nil {
			return
		}
	}
	s.assetsCache.WriteToCache(cacheKey, nil, content, requires, extInfo)
	return
}

func (s *Sprocket) readAssetContent(assetPath string, extInfo *types.ExtensionInfo) ([]byte, error) {
	content, err := ioutil.ReadFile(assetPath)
	if err != nil {
		return nil, errCantFindAsset(assetPath)
	}
	for _, f := range extInfo.ContentTreatment {
		content, err = f.Process(content, assetPath)
		if err != nil {
			return nil, err
		}
	}
	return content, nil
}
