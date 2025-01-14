# foggo

[![ci](https://github.com/ervitis/foggo/actions/workflows/ci.yml/badge.svg)](https://github.com/ervitis/foggo/actions/workflows/ci.yml)
[![Release](https://github.com/ervitis/foggo/actions/workflows/release.yml/badge.svg)](https://github.com/ervitis/foggo/actions/workflows/release.yml)
[![Coverage Status](https://coveralls.io/repos/github/ervitis/foggo/badge.svg?branch=main)](https://coveralls.io/github/ervitis/foggo?branch=main)

[日本語版 README](./README.ja.md)

__foggo__ is generator of `Functional Option Pattern` And `Applicable Functional Option Pattern` from struct field in Golang code.

## Installation

```shell
$ go install golang.org/x/tools/cmd/goimports@latest  # foggo use 'goimports' command
$ go install github.com/ervitis/foggo@latest
```

## Usage

__foggo__ provides `fop` and `afop` subcommand.

```shell
Usage:
  foggo (fop|afop) [flags]

Flags:
  -h, --help   help for fop

Global Flags:
  -p, --package string   Package name having target struct (default ".")
  -s, --struct string    Target struct name (required)
  -n, --no-instance bool Do not create the New method (default false)
```

### Generate with command line

1. prepare a struct type.

    ```go
    // ./image/image.go
    package image
    
    type Image struct {
        Width  int
        Height int
        // don't want to create option, specify `foggo:"-"` as the structure tag 
        Src    string `foggo:"-"`
        Alt    string
    }
    ```

2. execute `foggo fop` command.

    ```shell
    # struct must be set struct type name 
    # package must be package path
    $ foggo fop --struct Image --package image
    ```

3. then `foggo` generates Functional Option Pattern code to `./image/image_gen.go`.

    ```go
    // Code generated by foggo; DO NOT EDIT.

    package image

    type ImageOption func(*Image)

    func NewImage(options ...ImageOption) *Image {
        s := &Image{}
    
        for _, option := range options {
            option(s)
        }
    
        return s
    }
    
    func WithWidth(Width int) ImageOption {
        return func(args *Image) {
            args.Width = Width
        }
    }
    
    func WithHeight(Height int) ImageOption {
        return func(args *Image) {
            args.Height = Height
        }
    }

    func WithAlt(Alt string) ImageOption {
        return func(args *Image) {
            args.Alt = Alt
        }
    }
    ```
   
4. write Golang code using `functional option parameter`

    ```go
    package main
   
    import "github.com/user/project/image"
    
    func main() {
	    image := NewImage(
	    	WithWidth(1280),
	    	WithHeight(720),
	    	WithAlt("alt title"),
	    )
	    image.Src = "./image.png"
        ...
    }
    ```

### Generate with `go:generate`

1. prepare a struct type with `go:generate`.

    ```go
    // ./image/image.go
    package image
    
    //go:generate foggo fop --struct Image 
    type Image struct {
        Width  int
        Height int
        // don't want to create option, specify `foggo:"-"` as the structure tag 
        Src    string `foggo:"-"`
        Alt    string
    }
    ```

2. execute `go generate ./...` command.

    ```shell
    $ go generate ./...
    ```

3. the `foggo` generate Functional Option Pattern code to all files written `go:generate`. 

### Generate with `afop` command

`afop` is the method to generate `Applicable Functional Option Pattern` code.

1. prepare a struct type with `go:generate`. (use `afop` subcommand)

    ```go
    // ./image/image.go
    package image
    
    //go:generate foggo afop --struct Image 
    type Image struct {
        Width  int
        Height int
        // don't want to create option, specify `foggo:"-"` as the structure tag 
        Src    string `foggo:"-"`
        Alt    string
    }
    ```

2. execute `go generate ./...` command.

    ```shell
    $ go generate ./...
    ```

3. the `foggo` generate Applicable Functional Option Pattern code to all files written `go:generate`. 

    ```go
    // Code generated by foggo; DO NOT EDIT.

    package image

    type ImageOption interface {
        apply(*Image)
    }

    type WidthOption struct {
        Width int
    }

    func (o WidthOption) apply(s *Image) {
        s.Width = o.Width
    } 

    type HeightOption struct {
        Height int
    }

    func (o HeightOption) apply(s *Image) {
        s.Height = o.Height
    }

    type AltOption struct {
        Alt string
    }

    func (o AltOption) apply(s *Image) {
        s.Alt = o.Alt
    }

    func NewImage(options ...ImageOption) *Image {
        s := &Image{}

        for _, option := range options {
            option.apply(s)
        }

        return s
    }
    ```

4. write Golang code using `Applicable Functional Option Parameter`

    ```go
    package main
   
    import "github.com/user/project/image"
    
    func main() {
	    image := NewImage(
	    	WidthOption(1280),
	    	HeightOption(720),
	    	AltOption("alt title"),
	    )
	    image.Src = "./image.png"
        ...
    }
    ```


## Functional Option Pattern ?
`Functional Option Pattern`(`FOP`) is one of the most common design patterns used in Golang code.

Golang cannot provide optional arguments such as keyword arguments (available in python, ruby, ...).
`FOP` is the technique for achieving optional arguments.

For more information, please refer to the following articles.

- https://commandcenter.blogspot.jp/2014/01/self-referential-functions-and-design.html
- https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

### Applicable Functional Option Pattern ?

`Applicable Functional Option Pattern`(`AFOP`) is __testable__ `FOP`.
`FOP` express options to function.
For that reason, comparing to option function with same arguments fails (not testable).

`AFOP` express options to struct type and options have a parameter and `apply` method.
Struct type is comparable in Golang, options followed `AFOP` are testable. 

`AFOP` proposed by following articles.

- https://github.com/uber-go/guide/blob/master/style.md#functional-options
- https://ww24.jp/2019/07/go-option-pattern (in Japanese)

## References

- https://github.com/moznion/gonstructor
