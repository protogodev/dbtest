package builtin

type Testee struct {
	// Instance is a testee instance, which must implements the
	// corresponding interface.
	Instance interface{}
	Close    func() error
	NewDB    func() DB
	Codec    Codec
}
