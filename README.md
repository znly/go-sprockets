# <img src="logo.png" alt="logo" width="80px"> Go-Sprockets
[![GoDoc](https://godoc.org/github.com/znly/go-sprockets?status.png)](https://godoc.org/github.com/znly/go-sprockets)

Go-Sprockets is a Golang Library for for compiling and serving web assets.
It features declarative dependency management for JavaScript and CSS
assets, as well as a powerful preprocessor pipeline that allows you to
write assets in languages like CoffeeScript, Sass and SCSS.

It's a port of [rails/Sprockets Library](https://github.com/rails/sprockets).

## Installation
Get the package:

```bash
$ go get github.com/znly/go-sprockets
```

## How to use

An example equals thousand words!

Let s play with a [simplereader](examples/simplereader/main.go)
```go
package main

import (
    "fmt"
    "os"

    sprockets "github.com/znly/go-sprockets"
)

func main() {
    s, err := sprockets.NewWithDefault(os.Args[1], "")
    if err != nil {
        fmt.Println(err)
        return
    }
    content, err := s.GetAsset(os.Args[2])
    fmt.Println(string(content), err)
}
```

You can build it with ```go build -o simplereader examples/simplereader/main.go```

Look at the examples assets folder and subfolders and test those command lines
```bash
    ./simplereader examples/assets javascripts/coffeefile.coffee
    ./simplereader examples/assets coffeefile.coffee
    ./simplereader examples/assets coffeefile.js
    ./simplereader examples/assets outofjavascript.js
    ./simplereader examples/assets cssfile.css
    ./simplereader examples/assets afolder/infolder.js
```

Enjoy

## Compilation Pipeline
* Get the asset extension info
    * use the last \\..* as the extension
    * if no extension info is found, use the default one

* Search for the asset using the extenstion info
    * search for the asset in each path defined in the extension info
        * for each path try current extension then alternate extension
        * Stop and return the first result, changing the extension info to the new extension of the asset if needed
* Read the asset and all it s requirement

    * Read the raw asset
        * read the whole asset file content
        * apply all the content contentTreatment from the extension info

    * Find requirements
        * use the require pattern Head to retrieve the HEADER
        * apply all the header ContentTreatment from the extension info if HEADER is found
        * use the require pattern Requires to retrieve each directive lines in the HEADER
        * apply the filecompiler of the extension info
        * search for the requirements, read them and find their own requirements.

* Bundle and compile
    * resolve the dependency graph and build the full content of the asset
    * apply the bundlecompiler of the extension info
    * apply all the post compile ContentTreatment

* Return the result


## func NewWithDefault(assetsPath, publicPath string)
This function is here to mimic [rails/sprockets directive processor](https://github.com/rails/sprockets/blob/master/README.md#the-directive-processor) and you should read it

## Production Mode
Setting up a public path when creating a new SprocketGo will be use to return pre bundled assets.
If the asset is missing from the public path, it will be automatically build and saved in the public path

## Generate Public Assets
You can use the function ```func (*Sprocket) Generate(assetUrl string) (error)``` to force the generation of an asset from the asset path to the public path.
BEWARE: if public path is not set an error will be returned

## WARNING.
Go-Sprockets is using [go-libsass](http://github.com/wellington/go-libsass) which embeded a C library and thus may take some time to compile.  Use ```go install``` to avoid recompiling it too often.

Go-Sprockets is using [go-duktape](https://github.com/olebedev/go-duktape) which embeded a C library and thus may take some time to compile.  Use ```go install``` to avoid recompiling it too often.

## Todo:

* find extension info by iterating on each extension info name and check if the asset is ending by it (sort by size .min.js > .js)
* implement [Index files are proxies for folders](https://github.com/rails/sprockets#index-files-are-proxies-for-folders)
* implement [require self](https://github.com/rails/sprockets#the-require_self-directive)
* implement [depend_on](https://github.com/rails/sprockets#the-depend_on-directive)
* implement [depend_on_asset](https://github.com/rails/sprockets#the-depend_on_asset-directive)
* implement [stub](https://github.com/rails/sprockets#the-stub-directive)
* Write Tests
* Make the cache faster!

## License
Licensed under the Apache License, Version 2.0.
