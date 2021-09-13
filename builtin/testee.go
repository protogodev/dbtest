package builtin

import (
	"errors"
	"time"
)

type Testee struct {
	// Instance is a testee instance, which must implements the
	// corresponding interface.
	Instance interface{}
	DB       DB
	Codec    Codec
}

func (t *Testee) Complete() *Testee {
	if t.Codec == nil {
		t.Codec = &DefaultCodec{TimeFormat: time.RFC3339}
	}
	return t
}

func (t *Testee) Validate() error {
	if t == nil {
		return errors.New("t is nil")
	}
	if t.Instance == nil {
		return errors.New("t.Instance is nil")
	}
	if t.DB == nil {
		return errors.New("t.DB is nil")
	}
	return nil
}
