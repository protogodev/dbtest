package builtin

import (
	"github.com/RussellLuo/structool"
)

type Codec interface {
	Decode(in interface{}, out interface{}) (err error)
	Encode(in interface{}) (out interface{}, err error)
}

func DefaultCodec(layout string) Codec {
	return structool.New().TagName("dbtest").
		DecodeHook(
			structool.DecodeStringToError,
			structool.DecodeStringToTime(layout),
			structool.DecodeStringToDuration,
		).
		EncodeHook(
			structool.EncodeErrorToString,
			structool.EncodeTimeToString(layout),
			structool.EncodeDurationToString,
		)
}
