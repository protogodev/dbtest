package builtin

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var (
	errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

type Codec interface {
	Decode(in map[string]interface{}, out interface{}) (err error)
	Encode(in interface{}) (out map[string]interface{}, err error)
}

type DefaultCodec struct {
	decoder *mapstructure.Decoder
}

func (dc *DefaultCodec) Decode(in map[string]interface{}, out interface{}) error {
	return decode(in, out)
}

func (dc *DefaultCodec) Encode(in interface{}) (map[string]interface{}, error) {
	inValue := reflect.ValueOf(in)
	if inValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("non-struct is unsupported")
	}

	out := make(map[string]interface{})
	if err := decode(in, &out); err != nil {
		return nil, err
	}

	overwriteErr(inValue, out)
	return out, nil
}

func decode(in, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		//DecodeHook: fn,
		TagName: "dbtest",
		Result:  out,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(in)
}

// overwriteEr reset the possible error field from an error to a string.
func overwriteErr(in reflect.Value, out map[string]interface{}) {
	typ := in.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := in.Field(i)

		// NOTE: We assume that there is only one error at the top level.
		if field.Type.Implements(errorInterface) {
			name := getFieldName(field)
			var msg string
			if value.IsNil() {
				msg = ""
			} else {
				msg = value.Interface().(error).Error()
			}

			out[name] = msg
			return
		}
	}
}

func getFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("dbtest")
	if tag == "-" {
		return ""
	}

	parts := strings.Split(tag, ",")
	name := parts[0]

	if name == "" {
		name = field.Name
	}

	return name
}
