package libucl

import "errors"

// #include "go-libucl.h"
import "C"

// SchemaErrorCode is a progamatic way to ficgure out what went wrong during validation
type SchemaErrorCode int

const (
	// SchemaOK means nothing went wrong (no error)
	SchemaOK SchemaErrorCode = iota
	// SchemaTypeMismatch means the type of object is wrong
	SchemaTypeMismatch
	// SchemaInvalidSchema means the provided schema is not valid according to json-schema draft 4
	SchemaInvalidSchema
	// SchemaMissingProperty means at least one property of the object is missing
	SchemaMissingProperty
	// SchemaConstraint means a contraint was not met.
	SchemaConstraint
	// SchemaMissingDependency means a dependency was not met
	SchemaMissingDependency
	// SchemaGenericError is a generic error (matches UCL_SCHEMA_UNKNOWN in libucl)
	SchemaGenericError
)

// SchemaError contains information on an error found when validating an UCL Object
// against a provided json-schema style schema
type SchemaError struct {
	code    SchemaErrorCode
	message string
	object  *Object
}

// Validate validates the object againt a provided schema, which should conform to
// the 4th draft of the json-schema standard
func (o *Object) Validate(schema *Object) (SchemaError, error) {
	var cError C.ucl_schema_error_t
	var err error
	var schemaError SchemaError
	ok := C.ucl_object_validate(schema.object, o.object, &cError)
	if !ok {
		schemaError.code = SchemaErrorCode(cError.code)
		schemaError.message = bufferToString(cError.msg)
		schemaError.object = &Object{object: cError.obj}
		err = errors.New(schemaError.message)
	}
	return schemaError, err
}

// Convert a fixed char array to a go string, as gGo doesn't let us use
// [128]C.char as *C.char in C.GoString
func bufferToString(buffer [128]C.char) string {
	var byteBuffer []byte
	for _, c := range buffer {
		if c == 0 {
			break
		} else {
			byteBuffer = append(byteBuffer, byte(c))
		}
	}
	return string(byteBuffer)
}

// Purely exists for the test, as cgo isn't supported inside tests
func byteToBufferAdapter(buffer [128]byte) string {
	var tmp [128]C.char
	for i, c := range buffer {
		tmp[i] = C.char(c)
	}
	return bufferToString(tmp)
}
