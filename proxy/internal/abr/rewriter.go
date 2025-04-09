package abr

type ABRPolicy interface {
	RewriteMPD(original []byte, clientID string) ([]byte, error)
}

// dummy passthrough ARB policy
type PassthroughPolicy struct{}

func (p *PassthroughPolicy) RewriteMPD(original []byte, clientID string) ([]byte, error) {
	// No actual changes — passthrough
	return original, nil
}
