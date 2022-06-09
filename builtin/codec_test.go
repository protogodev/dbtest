package builtin_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/protogodev/dbtest/builtin"
)

func TestDefaultCodec_Decode(t *testing.T) {
	type Datum struct {
		Bool    bool
		Null    *int
		TimePtr *time.Time
	}
	type Out struct {
		Int   int
		Str   string
		Datum Datum
		Data  []Datum
		Time  time.Time
	}

	in := map[string]interface{}{
		"Int": 1,
		"Str": "s",
		"Datum": map[string]interface{}{
			"Bool":    true,
			"Null":    (*int)(nil),
			"TimePtr": "2021-08-13T00:00:00Z",
		},
		"Data": []interface{}{
			map[string]interface{}{
				"Bool":    true,
				"Null":    (*int)(nil),
				"TimePtr": "2021-08-13T00:00:00Z",
			},
		},
		"Time": "2021-08-13T00:00:00Z",
	}
	date := time.Date(2021, 8, 13, 0, 0, 0, 0, time.UTC)
	want := Out{
		Int:   1,
		Str:   "s",
		Datum: Datum{Bool: true, Null: nil, TimePtr: &date},
		Data:  []Datum{{Bool: true, Null: nil, TimePtr: &date}},
		Time:  date,
	}

	codec := builtin.DefaultCodec(time.RFC3339)
	var got Out
	if err := codec.Decode(in, &got); err != nil {
		t.Fatalf("Err: %v\n", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Out: Got (%#v) != Want (%#v)", got, want)
	}
}

func TestDefaultCodec_Encode(t *testing.T) {
	type Datum struct {
		Bool    bool
		Null    *int
		TimePtr *time.Time
	}
	type In struct {
		Int   int
		Str   string
		Datum Datum
		Data  []Datum
		Time  time.Time
		Err   error `dbtest:"err"`
	}

	date := time.Date(2021, 8, 13, 0, 0, 0, 0, time.UTC)
	in := In{
		Int:   1,
		Str:   "s",
		Datum: Datum{Bool: true, Null: nil, TimePtr: &date},
		Data:  []Datum{{Bool: true, Null: nil, TimePtr: nil}},
		Time:  date,
		Err:   errors.New("oops"),
	}
	want := map[string]interface{}{
		"Int": 1,
		"Str": "s",
		"Datum": map[string]interface{}{
			"Bool":    true,
			"Null":    (*int)(nil),
			"TimePtr": "2021-08-13T00:00:00Z",
		},
		"Data": []interface{}{
			map[string]interface{}{
				"Bool":    true,
				"Null":    (*int)(nil),
				"TimePtr": "",
			},
		},
		"Time": "2021-08-13T00:00:00Z",
		"err":  "oops",
	}

	codec := builtin.DefaultCodec(time.RFC3339)
	got, err := codec.Encode(in)
	if err != nil {
		t.Fatalf("Err: %v\n", err)
	}

	if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
		t.Fatalf("Out: Got (%#v) != Want (%#v)", got, want)
	}
}
