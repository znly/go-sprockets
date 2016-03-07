package sprockets

import (
	"path/filepath"

	"github.com/znly/go-sprockets/stringlist"
	"github.com/znly/go-sprockets/types"
)

func newExtensionInfo(ext string) *types.ExtensionInfo {
	ret := &types.ExtensionInfo{
		CurrentExtension: ext,
	}
	ret.Paths = stringlist.NewList()
	ret.AlterExts = stringlist.NewList()
	ret.AlterExts.PushFront(ext)
	return ret
}

func (s *Sprocket) getOrCreateExtensionInfo(ext string) *types.ExtensionInfo {
	if extInfo, ok := s.extInfos[ext]; ok {
		return extInfo
	}
	extInfo := newExtensionInfo(ext)
	extInfo.Paths.PushFrontList(s.defaultExtInfo.Paths)
	s.extInfos[ext] = extInfo
	return extInfo
}

func (s *Sprocket) getExtensionInfoOrDefault(ext string) *types.ExtensionInfo {
	if extInfo, ok := s.extInfos[ext]; ok {
		return extInfo
	}
	return s.defaultExtInfo
}

// PushFrontDefaultPath will add a path to the beginning of the list of default paths
// this path will be uniq in that list (old duplicate will be removed)
func (s *Sprocket) PushFrontDefaultPath(path string) (err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}
	s.defaultExtInfo.Paths.PushFrontUniq(path)
	return
}

// PushBackDefaultPath will add a path to the end of the list of default paths
// this path will be uniq in that list (old duplicate will be removed)
func (s *Sprocket) PushBackDefaultPath(path string) (err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}
	s.defaultExtInfo.Paths.PushBackUniq(path)
	return
}

// PushFrontExtensionPath will add a path to the beginning of the list of path for one Extension
// this path will be uniq in that list (old duplicate will be removed)
func (s *Sprocket) PushFrontExtensionPath(ext, path string) (err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}
	s.getOrCreateExtensionInfo(ext).Paths.PushFrontUniq(path)
	return
}

// PushBackExtensionPath will add a path to the end of the list of path for one Extension
// this path will be uniq in that list (old duplicate will be removed)
func (s *Sprocket) PushBackExtensionPath(ext, path string) (err error) {
	if path, err = filepath.Abs(path); err != nil {
		return
	}
	s.getOrCreateExtensionInfo(ext).Paths.PushBackUniq(path)
	return
}

// PushFrontAlterExtension will add a path to the beginning of the list of path for one Extension
// this path will be uniq in that list (old duplicate will be removed)
func (s *Sprocket) PushFrontAlterExtension(ext, alterExt string) {
	s.getOrCreateExtensionInfo(ext).AlterExts.PushFrontUniq(alterExt)
}

// PushBackAlterExtension will add an alter Extension to the end of the list of alter Extension for one Extension
// this alter Extension will be uniq in that list (old duplicate will be removed)
func (s *Sprocket) PushBackAlterExtension(ext, alterExt string) {
	s.getOrCreateExtensionInfo(ext).AlterExts.PushBackUniq(alterExt)
}

// AddContentTreatment will add a Content Treatment interface that will be use on the raw content of a file
func (s *Sprocket) AddContentTreatment(ext string, ContentTreatment types.ContentTreatmentInterface) {
	extInfo := s.getOrCreateExtensionInfo(ext)
	extInfo.ContentTreatment = append(extInfo.ContentTreatment, ContentTreatment)
}

// AddHeaderTreatment will add a Content Treatment interface that will be use on the require header of a file
func (s *Sprocket) AddHeaderTreatment(ext string, HeaderTreatment types.ContentTreatmentInterface) {
	extInfo := s.getOrCreateExtensionInfo(ext)
	extInfo.HeaderTreatment = append(extInfo.HeaderTreatment, HeaderTreatment)
}

// AddPostCompileContentTreatment will add a Content Treatment interface that will be use on content of a file after bundle Compiler
func (s *Sprocket) AddPostCompileContentTreatment(ext string, PostCompileContentTreatment types.ContentTreatmentInterface) {
	extInfo := s.getOrCreateExtensionInfo(ext)
	extInfo.PostCompileContentTreatment = append(extInfo.PostCompileContentTreatment, PostCompileContentTreatment)
}

// SetRequirePattern set the require pattern of this extension
func (s *Sprocket) SetRequirePattern(ext string, rp *types.RequirePattern) {
	extInfo := s.getOrCreateExtensionInfo(ext)
	extInfo.RequirePattern = rp
}

// SetBundleCompiler will add a Content Treatment interface that will be use on the full content of a file (means the content of this file and all it s requirements)
func (s *Sprocket) SetBundleCompiler(ext string, bundleCompiler types.ContentTreatmentInterface) {
	extInfo := s.getOrCreateExtensionInfo(ext)
	extInfo.BundleCompiler = bundleCompiler
}

// SetFileCompiler will add a Content Treatment interface that will be use on the content of a file after Content and Header Treatment
func (s *Sprocket) SetFileCompiler(ext string, fileCompiler types.ContentTreatmentInterface) {
	extInfo := s.getOrCreateExtensionInfo(ext)
	extInfo.FileCompiler = fileCompiler
}
