package sprockets

import (
	"path/filepath"
	"regexp"

	"github.com/znly/go-sprockets/bundlecompiler"
	"github.com/znly/go-sprockets/filecompiler"
	"github.com/znly/go-sprockets/types"
)

//NewWithDefault create a new Sprocket pipeline:
//- using assetsPath as default asset directory
//- ".css", ".scss" and ".sass" configuration
//    - adding filecompiler for sass (to turn it into scss)
//    - adding bundlecompiler for sass, scss, css
//    - adding a search path to [assetsPath]/stylesheets
//    - adding require rules
//- ".js" and ".coffee" configuration
//    - adding filecompiler for coffee script (to turn file into )
//    - adding require rules
//    - adding a search path to [assetsPath]/javascripts
//- ".jpg", ".png", ".svg", ".gif", ".bmp", ".tiff", ".tga" configuration
//    - adding a search path to [assetsPath]/images
//- ".eot", ".svg", ".ttf", ".woff" configuration
//    - adding a search path to [assetsPath]/fonts
func NewWithDefault(assetsPath, publicPath string) (s *Sprocket, err error) {
	s, err = New(assetsPath, publicPath)
	if err != nil {
		return
	}
	s.PushFrontAlterExtension(".css", ".sass")
	s.PushFrontAlterExtension(".css", ".scss")
	s.SetRequirePattern(".css", &types.RequirePattern{
		Head:    regexp.MustCompile(`(\s*(/\*(.*?\s*?)*\*/)*)*`),
		Require: regexp.MustCompile(`^\s*(?:\*/)?\s*(?:/\*.*?\*/)*\s*(\*\s*=\s*require((?:_directory|_tree)?)\s+(.+))`),
	})
	s.SetBundleCompiler(".css", &bundlecompiler.ScssSassCompiler{})
	s.PushFrontExtensionPath(".css", filepath.Join(s.assetsPath, "stylesheets"))

	s.PushFrontAlterExtension(".scss", ".css")
	s.PushFrontAlterExtension(".scss", ".sass")
	s.SetRequirePattern(".scss", &types.RequirePattern{
		Head:    regexp.MustCompile(`(\s*(/\*(.*\s+)*\*/)*([ \t]*//.*\s+)*)*`),
		Require: regexp.MustCompile(`^\s*(?:\*/)?\s*(?:/\*.*?\*/)*\s*((?:\*|//)\s*=\s*require((?:_directory|_tree)?)\s+(.+))`),
	})
	s.SetBundleCompiler(".scss", &bundlecompiler.ScssSassCompiler{})
	s.PushFrontExtensionPath(".scss", filepath.Join(s.assetsPath, "stylesheets"))

	s.PushFrontAlterExtension(".sass", ".css")
	s.PushFrontAlterExtension(".sass", ".scss")
	s.SetRequirePattern(".sass", &types.RequirePattern{
		Head:    regexp.MustCompile(`(\s*(/\*(.*\s+)*\*/)*([ \t]*//.*\s+)*)*`),
		Require: regexp.MustCompile(`^\s*(?:\*/)?\s*(?:/\*.*?\*/)*\s*((?:\*|//)\s*=\s*require((?:_directory|_tree)?)\s+(.+))`),
	})
	s.SetFileCompiler(".sass", &filecompiler.SassCompiler{})
	s.SetBundleCompiler(".sass", &bundlecompiler.ScssSassCompiler{})
	s.PushFrontExtensionPath(".sass", filepath.Join(s.assetsPath, "stylesheets"))

	s.PushFrontAlterExtension(".coffee", ".js")
	s.SetFileCompiler(".coffee", filecompiler.NewCoffeeCompiler())
	s.SetRequirePattern(".coffee", &types.RequirePattern{
		Head:    regexp.MustCompile(`(\s*#[^\n]*\n)*`),
		Require: regexp.MustCompile(`^(\s*#\s*=\s*require((?:_directory|_tree)?)\s+(.+))`),
	})
	s.PushFrontExtensionPath(".coffee", filepath.Join(s.assetsPath, "javascripts"))

	s.PushFrontAlterExtension(".js", ".coffee")
	s.SetRequirePattern(".js", &types.RequirePattern{
		Head:    regexp.MustCompile(`(\s*(/\*(.*\s+)*\*/)*([ \t]*//.*\s+)*)*`),
		Require: regexp.MustCompile(`^\s*(?:\*/)?\s*(?:/\*.*?\*/)*\s*((?:\*|//)\s*=\s*require((?:_directory|_tree)?)\s+(.+))`),
	})
	s.PushFrontExtensionPath(".js", filepath.Join(s.assetsPath, "javascripts"))

	for _, ext := range []string{".jpg", ".png", ".svg", ".gif", ".bmp", ".tiff", ".tga"} {
		s.PushFrontExtensionPath(ext, filepath.Join(s.assetsPath, "images"))
	}

	for _, ext := range []string{".eot", ".svg", ".ttf", ".woff"} {
		s.PushFrontExtensionPath(ext, filepath.Join(s.assetsPath, "fonts"))
	}
	return
}
