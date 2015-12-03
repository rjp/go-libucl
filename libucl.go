package libucl

// #cgo CFLAGS: -Ivendor/libucl/include -Wno-int-to-void-pointer-cast
// #cgo LDFLAGS: -Lvendor/libucl -lucl
//
// #cgo freebsd CFLAGS: -I/usr/local/include -Wno-int-to-void-pointer-cast
// #cgo freebsd LDFLAGS: -L/usr/local/lib -lucl
//
import "C"
