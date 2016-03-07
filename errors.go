package sprockets

import "errors"

//
func errCantFindAsset(path string) (err error) {
	err = errors.New("Cant find asset: " + path)
	return
}
