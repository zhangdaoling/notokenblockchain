package curver

import (
	"bytes"
	"errors"
	"math/rand"

	"golang.org/x/crypto/ed25519"
)

// errors
var (
	ErrWrongLengthEd25519PrivateKey = errors.New("Ed25519 private length is not 64")
	ErrInvalidEd25519PrivateKey     = errors.New("invalid Ed25519 privateKey")
	ErrGeneratePublicKey            = errors.New("generate Ed25519 publicKey error")
)

type Ed25519 struct{}

func (b *Ed25519) GeneratePrivateKey() ([]byte, error) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, nil
	}
	return ed25519.NewKeyFromSeed(seed), nil
}

func (b *Ed25519) CheckPrivateKey(privateKey []byte) error {
	if len(privateKey) != 64 {
		return ErrWrongLengthEd25519PrivateKey
	}
	if !bytes.Equal(ed25519.NewKeyFromSeed(privateKey[:32]), privateKey) {
		return ErrInvalidEd25519PrivateKey
	}
	return nil
}

func (b *Ed25519) GeneratePublicKey(privateKey []byte) ([]byte, error) {
	pubkey, ok := ed25519.PrivateKey(privateKey).Public().(ed25519.PublicKey)
	if !ok {
		return nil, ErrGeneratePublicKey
	}
	return pubkey, nil
}

func (b *Ed25519) Sign(message []byte, privateKey []byte) ([]byte, error) {
	return ed25519.Sign(privateKey, message), nil
}

func (b *Ed25519) Verify(message []byte, publicKey []byte, sig []byte) bool {
	return ed25519.Verify(publicKey, message, sig)
}
