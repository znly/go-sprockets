package bundlecompiler

import (
	"bytes"

	libsass "github.com/wellington/go-libsass"
)

// ScssSassCompiler is here to compile scss bundled file (sass must be transformed into scss first)
// it s also here to show you how to make a bundlecompiler
type ScssSassCompiler struct {
	LineNumbers bool
	DebugInfo   bool
}

// Process to implement ContentTreatmentInterface
func (ssc *ScssSassCompiler) Process(content []byte, path string) ([]byte, error) {
	in := bytes.NewBuffer(content)
	out := &bytes.Buffer{}
	comp, err := libsass.New(out, in, libsass.OutputStyle(libsass.Style["nested"]), libsass.Comments(ssc.LineNumbers || ssc.DebugInfo))
	if err != nil {
		return nil, err
	}
	if err := comp.Run(); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(bytes.Replace(out.Bytes(), []byte("*filecompiler.SassCompiler"), make([]byte, 0), 1)), nil
}
