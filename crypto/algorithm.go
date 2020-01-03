package crypto

import (
	"errors"
	"fmt"

	"github.com/zhangdaoling/simplechain/crypto/curver"
)

var (
	ErrVerifySignaturePanic = errors.New("verify signature panic")
	ErrNotSupportAlgorithm  = errors.New("not support algorithm")
)

type AlgorithmI interface {
	GeneratePrivateKey() ([]byte, error)
	CheckPrivateKey(privateKey []byte) error
	GeneratePublicKey(privateKey []byte) ([]byte, error)
	Sign(message []byte, privateKey []byte) ([]byte, error)
	Verify(message []byte, publicKey []byte, sig []byte) bool
}

var _ AlgorithmI = &curver.Ed25519{}
var _ AlgorithmI = Ed25519

type Algorithm int32

const (
	_ Algorithm = iota
	Ed25519
)

func NewAlgorithm(name string) (Algorithm, error) {
	switch name {
	case "ed25519":
		return Ed25519, nil
	default:
		return 0, ErrNotSupportAlgorithm
	}
}

func (a Algorithm) getBackend() (AlgorithmI, error) {
	switch a {
	case Ed25519:
		return &curver.Ed25519{}, nil
	default:
		return &curver.Ed25519{}, ErrNotSupportAlgorithm
	}
}

func (a Algorithm) String() string {
	switch a {
	case Ed25519:
		return "ed25519"
	default:
		return ErrNotSupportAlgorithm.Error()
	}
}

func (a Algorithm) GeneratePrivateKey() ([]byte, error) {
	i, err := a.getBackend()
	if err != nil {
		return nil, err
	}
	return i.GeneratePrivateKey()
}

func (a Algorithm) CheckPrivateKey(privateKey []byte) error {
	i, err := a.getBackend()
	if err != nil {
		return err
	}
	return i.CheckPrivateKey(privateKey)
}

// GetPubkey will get the public key of the secret key
func (a Algorithm) GeneratePublicKey(privateKey []byte) ([]byte, error) {
	i, err := a.getBackend()
	if err != nil {
		return nil, err
	}
	return i.GeneratePublicKey(privateKey)
}

func (a Algorithm) Sign(message []byte, privateKey []byte) ([]byte, error) {
	i, err := a.getBackend()
	if err != nil {
		return nil, err
	}
	return i.Sign(message, privateKey)
}

func (a Algorithm) Verify(message []byte, publicKey []byte, sig []byte) (ret bool) {
	i, err := a.getBackend()
	if err != nil {
		return false
	}
	// catch ed25519.Verify panic
	defer func() {
		if e := recover(); e != nil {
			fmt.Errorf("verify panic. err=%v", e)
			ret = false
		}
	}()
	return i.Verify(message, publicKey, sig)
}
