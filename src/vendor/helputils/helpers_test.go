package helputils

import (
	"strings"
	"testing"
)

func TestMustStructToKeyValueInterfaces(t *testing.T) {
	s := &struct {
		A int
		b int
		C string
		d string
	}{
		A: 1,
		b: 2,
		C: "c",
		d: "d",
	}
	t.Log("struct:", s)
	kvIfs := MustStructToKeyValueInterfaces(s)
	if len(kvIfs) != 4 {
		t.Fatalf("unexpected kv len: %d", len(kvIfs))
	}
	t.Log("kv ifs:", len(kvIfs), kvIfs)

	defer func() {
		if r := recover(); r != nil {
			if strings.Contains(r.(string), "must have kind struct") {
				t.Log("expected failure:", r)
			} else {
				t.Fatal(r)
			}
		}
	}()
	MustStructToKeyValueInterfaces(2)
}
