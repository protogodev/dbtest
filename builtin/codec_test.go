package builtin_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/RussellLuo/dbtest/builtin"
)

func TestDefaultCodec_Encode(t *testing.T) {
	type Data struct {
		Bool bool
		Null *int
	}
	type in struct {
		Int  int
		Str  string
		Data Data
		Err  error `dbtest:"err"`
	}

	codec := &builtin.DefaultCodec{}
	got, err := codec.Encode(in{
		Int:  1,
		Str:  "s",
		Data: Data{Bool: true, Null: nil},
		Err:  errors.New("oops"),
	})
	if err != nil {
		t.Fatalf("Err: %v\n", err)
	}

	want := map[string]interface{}{
		"Int": 1,
		"Str": "s",
		"Data": map[string]interface{}{
			"Bool": true,
			"Null": (*int)(nil),
		},
		"err": "oops",
	}
	if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
		t.Fatalf("Out: Got (%#v) != Want (%#v)", got, want)
	}
}
