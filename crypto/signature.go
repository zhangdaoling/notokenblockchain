package crypto

import ()

// Signature is the signature of some message
type Signature struct {
	Algorithm Algorithm

	Sig    []byte
	Pubkey []byte
}
