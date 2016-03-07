package types

import "regexp"

// RequirePattern is a structure with regexp to find the sprocket's directive.
type RequirePattern struct {
	Head    *regexp.Regexp
	Require *regexp.Regexp
}

// RequireInterface need to be implemented to return the list of files of a sprocket's directive line.
type RequireInterface interface {
	//TODO: maybe use a list here too :) (when the merge list will work perfectly)
	//Return a list of path sorted and the timestamp of the last modified file or directory.
	GetList(*ExtensionInfo) ([]string, int64, error)
}
