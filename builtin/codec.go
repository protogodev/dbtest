package builtin

import (
	"fmt"
	"reflect"
	"time"

	"github.com/RussellLuo/structs"
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
	TimeFormat string
}

func (dc *DefaultCodec) Decode(in map[string]interface{}, out interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: dc.stringToTimeHookFunc,
		TagName:    "dbtest",
		Result:     out,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(in)
}

func (dc *DefaultCodec) Encode(in interface{}) (map[string]interface{}, error) {
	inValue := reflect.ValueOf(in)
	if inValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("non-struct is unsupported")
	}

	s := structs.New(in)
	s.TagName = "dbtest"
	s.EncodeHook = dc.errOrTimeToStringHookFunc

	return s.Map(), nil
}

func (dc *DefaultCodec) stringToTimeHookFunc(from, to reflect.Value) (interface{}, error) {
	if from.Kind() != reflect.String {
		return from.Interface(), nil
	}

	value := from.Interface().(string)

	switch to.Interface().(type) {
	case time.Time:
		return time.Parse(dc.TimeFormat, value)
	case *time.Time:
		t, err := time.Parse(dc.TimeFormat, value)
		if err != nil {
			return nil, err
		}
		return &t, nil
	}

	return from.Interface(), nil
}

func (dc *DefaultCodec) errOrTimeToStringHookFunc(value reflect.Value) (interface{}, error) {
	in := value.Interface()

	if value.Type().Implements(errorInterface) {
		if in == nil {
			return "", nil
		}
		return in.(error).Error(), nil
	}

	switch v := in.(type) {
	case time.Time:
		return v.Format(dc.TimeFormat), nil
	case *time.Time:
		if v == nil {
			return "", nil
		}
		return v.Format(dc.TimeFormat), nil
	}
	return in, nil
}
