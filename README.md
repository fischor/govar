# govar

A simple tool to statically print out values of consts and vars defined in generic declaration nodes of Go code.
No build of the files to operate on is required.

## Usage

```
govar FILE VARNAME
```

## Example

```golang
// A file named version/info.go
package version

var (
        myVersion = "1.0.0"
)
```

```sh
$ govar version/info.go myVersion
1.0.0
```

## Install

```
go install github.com/fischor/govar@latest
```

## Notes

Variables must be defined within a `var` or `const` block or line and must be of an assignment of a numeric value, a char or a string.
References to other variables will not resolved and would cause an error.
