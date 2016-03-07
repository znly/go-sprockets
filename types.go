package sprockets

import (
	"github.com/znly/go-sprockets/assetscache"
	"github.com/znly/go-sprockets/types"
)

// Sprocket structure use in sprocketgo
type Sprocket struct {
	assetPath      string
	extInfos       map[string]*types.ExtensionInfo
	defaultExtInfo *types.ExtensionInfo
	publicPath     string
	assetsCache    *assetscache.AssetsCache
}
