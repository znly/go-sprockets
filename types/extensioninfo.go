package types

import "github.com/znly/go-sprockets/stringlist"

// ContentTreatmentInterface is an interface needed for Content Treatment while processing files.
type ContentTreatmentInterface interface {
	//Process
	// will do the treatment on the content.
	// For BundleCompiler, path will be the root file path
	Process(content []byte, path string) ([]byte, error)
}

// ExtensionInfo is the configuration structure to know how sprocketgo need to read/compile an asset base on its extension
type ExtensionInfo struct {
	CurrentExtension            string
	Paths                       *stringlist.List
	RequirePattern              *RequirePattern
	AlterExts                   *stringlist.List
	ContentTreatment            []ContentTreatmentInterface
	HeaderTreatment             []ContentTreatmentInterface
	PostCompileContentTreatment []ContentTreatmentInterface
	BundleCompiler              ContentTreatmentInterface
	FileCompiler                ContentTreatmentInterface
}
