package libucl

import (
	"io/ioutil"
	"path"
	"testing"
)

func testParseString(t *testing.T, data string) *Object {
	obj, err := ParseString(data)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return obj
}

func TestParser(t *testing.T) {
	p := NewParser(0)
	defer p.Close()

	if err := p.AddString(`foo = bar;`); err != nil {
		t.Fatalf("err: %s", err)
	}

	obj := p.Object()
	if obj == nil {
		t.Fatal("obj should not be nil")
	}
	defer obj.Close()

	if obj.Type() != ObjectTypeObject {
		t.Fatalf("bad: %#v", obj.Type())
	}

	value := obj.Get("foo")
	if value == nil {
		t.Fatal("should have value")
	}
	defer value.Close()

	if value.Type() != ObjectTypeString {
		t.Fatalf("bad: %#v", obj.Type())
	}

	if value.Key() != "foo" {
		t.Fatalf("bad: %#v", value.Key())
	}

	if value.ToString() != "bar" {
		t.Fatalf("bad: %#v", value.ToString())
	}
}

func TestParserAddFile(t *testing.T) {
	tf, err := ioutil.TempFile("", "libucl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tf.Write([]byte("foo = bar;"))
	tf.Close()

	p := NewParser(0)
	defer p.Close()

	if err := p.AddFile(tf.Name()); err != nil {
		t.Fatalf("err: %s", err)
	}

	obj := p.Object()
	if obj == nil {
		t.Fatal("obj should not be nil")
	}
	defer obj.Close()

	if obj.Type() != ObjectTypeObject {
		t.Fatalf("bad: %#v", obj.Type())
	}

	value := obj.Get("foo")
	if value == nil {
		t.Fatal("should have value")
	}
	defer value.Close()

	if value.Type() != ObjectTypeString {
		t.Fatalf("bad: %#v", obj.Type())
	}

	if value.Key() != "foo" {
		t.Fatalf("bad: %#v", value.Key())
	}

	if value.ToString() != "bar" {
		t.Fatalf("bad: %#v", value.ToString())
	}
}

func TestParserRegisterMacro(t *testing.T) {
	value := ""
	parameter := ""
	macro := func(args Object, body string) bool {
		thing := args.Get("thing")
		if nil == thing {
			return false
		}
		parameter = thing.ToString()
		value = body
		return true
	}

	config := `.foo(thing=something) "bar";`

	p := NewParser(0)
	defer p.Close()

	p.RegisterMacro("foo", macro)

	if err := p.AddString(config); err != nil {
		t.Fatalf("err: %s", err)
	}

	if value != "bar" {
		t.Fatalf("body bad: %#v", value)
	}
	if parameter != "something" {
		t.Fatalf("parameter bad: %#v", parameter)
	}
}

func TestParseString(t *testing.T) {
	obj, err := ParseString("foo = bar; baz = boo;")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if obj == nil {
		t.Fatal("should have object")
	}
	defer obj.Close()

	if obj.Len() != 2 {
		t.Fatalf("bad: %d", obj.Len())
	}
}

func TestRegisterVariable(t *testing.T) {
	p := NewParser(0)
	defer p.Close()
	value := "bar"
	p.RegisterVariable("FOO", value)

	err := p.AddString("foo = $FOO; baz = boo;")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	obj := p.Object()
	if obj == nil {
		t.Fatal("Configuration should produce object")
	}
	defer obj.Close()

	v := obj.Get("foo")
	if v == nil {
		t.Fatal("Key \"foo\" should exist")
	}
	defer v.Close()
	if v.ToString() != value {
		t.Fatalf("bad: \"%s\", expected: \"%s\"", v.ToString(), value)
	}
}

func TestFileVariables(t *testing.T) {
	tf, err := ioutil.TempFile("", "libucl")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tf.Write([]byte("file = $FILENAME; dir = $CURDIR"))
	tf.Close()

	p := NewParser(0)
	defer p.Close()

	if err := p.AddFileAndSetVariables(tf.Name(), false); err != nil {
		t.Fatalf("err: %s", err)
	}

	obj := p.Object()
	if obj == nil {
		t.Fatal("obj should not be nil")
	}
	defer obj.Close()

	filename := obj.Get("file")
	if filename == nil {
		t.Fatal("key \"file\" should exist")
	}
	defer filename.Close()
	if filename.ToString() != tf.Name() {
		t.Errorf("bad: %s, expected %s", filename.ToString(), tf.Name())
	}

	dir := obj.Get("dir")
	if dir == nil {
		t.Fatal("key \"dir\" should exist")
	}
	defer dir.Close()
	if dir.ToString() != path.Dir(tf.Name()) {
		t.Errorf("bad: %s, expected %s", dir.ToString(), path.Dir(tf.Name()))
	}
}
