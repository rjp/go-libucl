// Package libucl provides golang bindings to libucl, a configuration library for
// UCL, the Universal Configuration Language.
package libucl

// #cgo CFLAGS: -Wno-int-to-void-pointer-cast
// #cgo pkg-config: libucl
import "C"
