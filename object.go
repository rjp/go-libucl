package libucl

import "unsafe"

// #include "go-libucl.h"
import "C"

// Object represents a single object within a configuration.
type Object struct {
	object *C.ucl_object_t
}

// ObjectIter is an interator for objects.
type ObjectIter struct {
	expand bool
	object *C.ucl_object_t
	iter   C.ucl_object_iter_t
}

// ObjectType is an enum of the type that an Object represents.
type ObjectType int

const (
	// ObjectTypeObject signifies a UCL Object (key/value pair)
	ObjectTypeObject ObjectType = iota
	// ObjectTypeArray signifies a UCL array
	ObjectTypeArray
	// ObjectTypeInt signifies an integer number
	ObjectTypeInt
	// ObjectTypeFloat signifies a floating-point nmber
	ObjectTypeFloat
	// ObjectTypeString signifies a string
	ObjectTypeString
	// ObjectTypeBoolean signifies a boolean value (true/false)
	ObjectTypeBoolean
	// ObjectTypeTime signifies time in seconds stored as a floating-point number
	ObjectTypeTime
	// ObjectTypeUserData signifies an opaque user-provided pointer, typically
	// used in macros
	ObjectTypeUserData
	// ObjectTypeNull signifies a null/non-existant value
	ObjectTypeNull
)

// Emitter is a type of built-in emitter that can be used to convert
// an object to another config format. All Emitters except EmitConfig
// are considered lossy, and information such as implicit arrays can be lost.
type Emitter int

const (
	// EmitJSON is the canonic json notation (with spaces indented structure)
	EmitJSON Emitter = iota
	// EmitJSONCompact is compact json notation (without spaces or newlines)
	EmitJSONCompact
	// EmitConfig is UCL (nginx-like)
	EmitConfig
	// EmitYAML is yaml inlined notation
	EmitYAML
)

// Close frees the memory associated with the object. This must be called when
// you're done using it.
func (o *Object) Close() error {
	C.ucl_object_unref(o.object)
	return nil
}

// Emit converts this object to another format and returns it.
func (o *Object) Emit(t Emitter) (string, error) {
	result := C.ucl_object_emit(o.object, uint32(t))
	if result == nil {
		return "", nil
	}

	return C.GoString(C._go_uchar_to_char(result)), nil
}

// Delete removes the given key from the object. The key will automatically
// be dereferenced once when this is called.
func (o *Object) Delete(key string) {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	C.ucl_object_delete_key(o.object, ckey)
}

// Get returns the element with matching key.
func (o *Object) Get(key string) *Object {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	obj := C.ucl_object_find_keyl(o.object, ckey, C.size_t(len(key)))
	if obj == nil {
		return nil
	}

	result := &Object{object: obj}
	result.Ref()
	return result
}

// Iterate returns an iterator that iterates over the objects in this object.
//
// The iterator must be closed when it is finished.
//
// The iterator does not need to be fully consumed.
func (o *Object) Iterate(expand bool) *ObjectIter {
	// Increase the ref count
	C.ucl_object_ref(o.object)

	return &ObjectIter{
		expand: expand,
		object: o.object,
		iter:   nil,
	}
}

// Key returns the key of this value/object as a string, or the empty
// string if the object doesn't have a key.
func (o *Object) Key() string {
	return C.GoString(C.ucl_object_key(o.object))
}

// Len returns the length of the object, or how many elements are part
// of this object.
//
// For objects, this is the number of key/value pairs.
// For arrays, this is the number of elements.
func (o *Object) Len() uint {
	// This is weird. If the object is an object and it has a "next",
	// then it is actually an array of objects, and to get the count
	// we actually need to iterate and count.
	if o.Type() == ObjectTypeObject && o.object.next != nil {
		iter := o.Iterate(false)
		defer iter.Close()

		var count uint
		for obj := iter.Next(); obj != nil; obj = iter.Next() {
			obj.Close()
			count++
		}

		return count
	}

	return uint(o.object.len)
}

// Ref increments the ref count associated with this. You have to call
// close an additional time to free the memory.
func (o *Object) Ref() error {
	C.ucl_object_ref(o.object)
	return nil
}

// Type returns the type that this object represents.
func (o *Object) Type() ObjectType {
	return ObjectType(C.ucl_object_type(o.object))
}

//------------------------------------------------------------------------
// Conversion Functions
//------------------------------------------------------------------------

// ToBool converts a UCL Object to a boolean value
func (o *Object) ToBool() bool {
	return bool(C.ucl_object_toboolean(o.object))
}

// ToInt converts a UCL Object to a signed integer value
func (o *Object) ToInt() int64 {
	return int64(C.ucl_object_toint(o.object))
}

// ToUint converts a UCL Object to an unsigned integer value
func (o *Object) ToUint() uint64 {
	return uint64(C.ucl_object_toint(o.object))
}

// ToFloat converts a UCL Object to an floating point value
func (o *Object) ToFloat() float64 {
	return float64(C.ucl_object_todouble(o.object))
}

// ToString converts a UCL Object to a string
func (o *Object) ToString() string {
	return C.GoString(C.ucl_object_tostring(o.object))
}

// Close frees the object iterator
func (o *ObjectIter) Close() {
	C.ucl_object_unref(o.object)
}

// Next returns the next iterative UCL Object
func (o *ObjectIter) Next() *Object {
	obj := C.ucl_object_iterate_with_error(o.object, &o.iter, C._Bool(o.expand), nil)
	if obj == nil {
		return nil
	}

	// Increase the ref count so we have to free it
	C.ucl_object_ref(obj)

	return &Object{object: obj}
}

// StringFlag are flags used in the conversion of strings into UCL objects
type StringFlag int

const (
	// StringEscape tells the converter to JSON escape the inputed string
	StringEscape StringFlag = C.UCL_STRING_ESCAPE
	// StringTrim tells the converter to trim leading and trailing whitespaces
	StringTrim StringFlag = C.UCL_STRING_TRIM
	// StringParseBoolean tells the converter to parse the inputted string as a boolean
	StringParseBoolean StringFlag = C.UCL_STRING_PARSE_BOOLEAN
	// StringParseInt tells the converter to parse the inputted string as an integer
	StringParseInt StringFlag = C.UCL_STRING_PARSE_INT
	// StringParseDouble tells the converter to parse the inputted string as a
	// floating-point number
	StringParseDouble StringFlag = C.UCL_STRING_PARSE_DOUBLE
	// StringParseTime tells the converter to parse the inputted string as a
	// time value, and treat as a floating-point number.
	StringParseTime StringFlag = C.UCL_STRING_PARSE_TIME
	// StringParseNumber tells the converter to parse the inputted string as a
	// number (integer, floating-point or time)
	StringParseNumber StringFlag = C.UCL_STRING_PARSE_TIME
	// StringParse tells the converter to parse the inputted string
	StringParse StringFlag = C.UCL_STRING_PARSE
	// StringParseBytes tells the converter to parse the inputted string as being
	// in bytes notation (e.g. 10k = 10*1024, not 10*1000)
	StringParseBytes StringFlag = C.UCL_STRING_PARSE_BYTES
)

// NewObject creates a new UCL Object from a string, JSON escaping it in the process
func NewObject(data string) *Object {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))
	obj := C.ucl_object_fromlstring(cData, C.size_t(len(data)))
	return &Object{object: obj}
}

// NewFormattedObject creates a new UCL Object from a string, according to the instructions
// given in the flags
func NewFormattedObject(data string, flags StringFlag) *Object {
	cData := C.CString(data)
	defer C.free(unsafe.Pointer(cData))
	obj := C.ucl_object_fromstring_common(cData, C.size_t(len(data)), uint32(flags))
	return &Object{object: obj}
}

// NewIntegerObject creates a new UCL Object from a 64-bit integer
func NewIntegerObject(data int64) *Object {
	obj := C.ucl_object_fromint(C.int64_t(data))
	return &Object{object: obj}
}

// NewDoubleObject creates a new UCL Object from a 64-bit floating-point number
func NewDoubleObject(data float64) *Object {
	obj := C.ucl_object_fromdouble(C.double(data))
	return &Object{object: obj}
}

// NewBoolObject creates a new UCL Object from a boolean
func NewBoolObject(data bool) *Object {
	obj := C.ucl_object_frombool(C.bool(data))
	return &Object{object: obj}
}
