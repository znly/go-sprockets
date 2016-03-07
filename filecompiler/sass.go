package filecompiler

import (
	"bytes"

	libsass "github.com/wellington/go-libsass"
)

// SassCompiler is here to compile a sass file into a scss file
// it s also here to show you how to make a file compiler
type SassCompiler struct {
}

// Process to implement ContentTreatmentInterface
func (sc *SassCompiler) Process(content []byte, path string) ([]byte, error) {
	in := bytes.NewBuffer(content)
	out := &bytes.Buffer{}
	if err := libsass.ToScss(in, out); err != nil {
		return nil, err
	}
	return bytes.TrimSpace(out.Bytes()), nil
}
