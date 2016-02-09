# Libucl Library for Go

go-libucl is a [libucl](https://github.com/vstakhov/libucl) library for
[Go](http://golang.org). Rather than re-implement libucl in Go, this library
uses cgo to bind directly to libucl. This allows the libucl project to be
the central source of knowledge. This project works on Mac OS X, Linux, and
Windows.

**Warning:** This library is still under development and API compatibility
is not guaranteed. Additionally, it is not feature complete yet, though
it is certainly usable for real purposes (we do!).

## Prerequisites
* libucl (This is a wrapper for this library)
* pkg-config (cgo uses this for locate where libucl is)

## Installation

```
$ go get github.com/draringi/go-libucl
```

Documentation is available on GoDoc: http://godoc.org/github.com/draringi/go-libucl
