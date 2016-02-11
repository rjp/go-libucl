# Libucl Library for Go
[![Build Status](https://travis-ci.org/draringi/go-libucl.svg?branch=master)](https://travis-ci.org/draringi/go-libucl)
[![GoDoc](https://godoc.org/github.com/draringi/go-libucl?status.svg)](https://godoc.org/github.com/draringi/go-libucl)

This version of go-libucl is forked from the [mitchellh version](https://github.com/mitchellh/go-libucl),
with the goal of having a version with a focus on using the OS's copy of libucl, in a portable manner,
as well as improve the Documentation quality.
As such, it uses pkg-config to determine the location of libucl.
In addition, patches from [hwx](https://github.com/bitmark-inc/go-libucl) have been
pulled in with 64-bit integer support.

go-libucl is a [libucl](https://github.com/vstakhov/libucl) library for
[Go](http://golang.org). Rather than re-implement libucl in Go, this library
uses cgo to bind directly to libucl. This allows the libucl project to be
the central source of knowledge. This project has been tested on Linux and FreeBSD.

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
