package libucl

import (
	"errors"
	"os"
	"sync"
	"unsafe"
)

// #include "go-libucl.h"
import "C"

// MacroFunc is the callback type for macros.
// return true os string is valid
// a macro call looks like:
//
//   .macro(UCL-OBJECT) "body-text"
//   .macro(key=value) "body-text"
//   .macro(params={key1=value1;key2=value2}) "body"
//
type MacroFunc func(args Object, body string) bool

// ParserFlag are flags that can be used to initialize a parser.
type ParserFlag int

const (
	// ParserKeyLowercase will lowercase all keys.
	ParserKeyLowercase ParserFlag = C.UCL_PARSER_KEY_LOWERCASE
	// ParserZeroCopy will attempt to do a zero-copy parse if possible.
	ParserZeroCopy ParserFlag = C.UCL_PARSER_ZEROCOPY
	// ParserNoTime will treat time values as strings.
	ParserNoTime ParserFlag = C.UCL_PARSER_NO_TIME
	// ParserNoImplicitArrays forces the creation explicit arrays instead of
	// implicit ones
	ParserNoImplicitArrays ParserFlag = C.UCL_PARSER_NO_IMPLICIT_ARRAYS
)

// Keeps track of all the macros internally
var macros map[int]MacroFunc
var macrosIdx int
var macrosLock sync.Mutex

// Parser is responsible for parsing libucl data.
type Parser struct {
	macros []int
	parser *C.struct_ucl_parser
}

// ParseString parses a string and returns the top-level object.
func ParseString(data string) (*Object, error) {
	p := NewParser(0)
	defer p.Close()
	if err := p.AddString(data); err != nil {
		return nil, err
	}

	return p.Object(), nil
}

// NewParser returns a parser
func NewParser(flags ParserFlag) *Parser {
	return &Parser{
		parser: C.ucl_parser_new(C.int(flags)),
	}
}

// AddString adds a string data to parse.
func (p *Parser) AddString(data string) error {
	cs := C.CString(data)
	defer C.free(unsafe.Pointer(cs))

	result := C.ucl_parser_add_string(p.parser, cs, C.size_t(len(data)))
	if !result {
		errstr := C.ucl_parser_get_error(p.parser)
		return errors.New(C.GoString(errstr))
	}
	return nil
}

// AddFile adds a file to parse.
func (p *Parser) AddFile(path string) error {
	cs := C.CString(path)
	defer C.free(unsafe.Pointer(cs))

	result := C.ucl_parser_add_file(p.parser, cs)
	if !result {
		errstr := C.ucl_parser_get_error(p.parser)
		return errors.New(C.GoString(errstr))
	}
	return nil
}

// Close frees the parser. Once it is freed it can no longer be used. You
// should always free the parser once you're done with it to clean up
// any unused memory.
func (p *Parser) Close() {
	C.ucl_parser_free(p.parser)

	if len(p.macros) > 0 {
		macrosLock.Lock()
		defer macrosLock.Unlock()
		for _, idx := range p.macros {
			delete(macros, idx)
		}
	}
}

// Object retrieves the root-level object for a configuration.
func (p *Parser) Object() *Object {
	obj := C.ucl_parser_get_object(p.parser)
	if obj == nil {
		return nil
	}

	return &Object{object: obj}
}

// RegisterMacro registers a macro that is called from the configuration.
func (p *Parser) RegisterMacro(name string, f MacroFunc) {
	// Register it globally
	macrosLock.Lock()
	if macros == nil {
		macros = make(map[int]MacroFunc)
	}
	for macros[macrosIdx] != nil {
		macrosIdx++
	}
	idx := macrosIdx
	macros[idx] = f
	macrosIdx++
	macrosLock.Unlock()

	// Register the index with our parser so we can free it
	p.macros = append(p.macros, idx)

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	C.ucl_parser_register_macro(
		p.parser,
		cname,
		C._go_macro_handler_func(),
		C._go_macro_index(C.int(idx)))
}

//export go_macro_call
func go_macro_call(id C.int, arguments *C.ucl_object_t, data *C.char, n C.int) C.bool {
	macrosLock.Lock()
	f := macros[int(id)]
	macrosLock.Unlock()

	args := Object{
		object: arguments,
	}

	// Macro not found, return error
	if f == nil {
		return false
	}

	// Macro found, call it!
	f(args, C.GoStringN(data, n))
	return true
}

// SetFileVariables sets the standard file variables ($FILENAME and $CURDIR) based
// on the provided filepath. If the argument expand is true, the path will be expanded
// out to an absolute path
//
// For example, if the current directory is /etc/nginx, and you give a path of
// ../file.conf, with exand = false, $FILENAME = ../file.conf and $CURDIR = ..,
// while with expand = true, $FILENAME = /etc/file.conf and $CURDIR = /etc
func (p *Parser) SetFileVariables(filepath string, expand bool) error {
	cpath := C.CString(filepath)
	defer C.free(unsafe.Pointer(cpath))
	result := C.ucl_parser_set_filevars(p.parser, cpath, C.bool(expand))
	if !result {
		errstr := C.ucl_parser_get_error(p.parser)
		return errors.New(C.GoString(errstr))
	}
	return nil
}

// RegisterVariable adds a new variable to the parser, which can be accessed in
// the configuration file as $variable_name
func (p *Parser) RegisterVariable(variable, value string) {
	cVariable := C.CString(variable)
	defer C.free(unsafe.Pointer(cVariable))
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	C.ucl_parser_register_variable(p.parser, cVariable, cValue)
}

// AddFileAndSetVariables is a combination of AddFile and SetFileVariables.
// It is meant to be a simple way to do both actions in a single function call.
func (p *Parser) AddFileAndSetVariables(path string, expand bool) error {
	err := p.AddFile(path)
	if err != nil {
		return err
	}

	err = p.SetFileVariables(path, expand)
	return err
}

// AddOpenFile reads in the configuration from a file already opened using os.Open
// or a related function.
func (p *Parser) AddOpenFile(f *os.File) error {
	fd := f.Fd()
	result := C.ucl_parser_add_fd(p.parser, C.int(fd))
	if !result {
		errstr := C.ucl_parser_get_error(p.parser)
		return errors.New(C.GoString(errstr))
	}
	return nil
}
