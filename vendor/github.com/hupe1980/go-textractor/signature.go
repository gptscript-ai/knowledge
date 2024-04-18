package textractor

// Compile time check to ensure Signature satisfies the LayoutChild interface.
var _ LayoutChild = (*Signature)(nil)

type Signature struct {
	base
}

func (s *Signature) Words() []*Word {
	return nil
}

func (s *Signature) Text(optFns ...func(*TextLinearizationOptions)) string {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	return opts.SignatureToken
}
