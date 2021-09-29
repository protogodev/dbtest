package builtin

import (
	"errors"
	"time"
)

type Testee struct {
	// System under test (http://xunitpatterns.com/SUT.html),
	// which must implement the corresponding interface.
	SUT   interface{}
	DB    DB
	Codec Codec
}

func (t *Testee) Complete() *Testee {
	if t.Codec == nil {
		t.Codec = DefaultCodec(time.RFC3339)
	}
	return t
}

func (t *Testee) Validate() error {
	if t == nil {
		return errors.New("t is nil")
	}
	if t.SUT == nil {
		return errors.New("t.SUT is nil")
	}
	if t.DB == nil {
		return errors.New("t.DB is nil")
	}
	return nil
}
