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
	GeneratePrivateKey() (error, []byte)
	CheckPrivateKey(privateKey []byte) error
	GeneratePublicKey(privateKey []byte) (error, []byte)
	Sign(message []byte, privateKey []byte) (error, []byte)
	Verify(message []byte, publicKey []byte, sig []byte) bool
}

var _ AlgorithmI = &curver.Ed25519{}
var _ AlgorithmI = Ed25519

type Algorithm uint8

const (
	_ Algorithm = iota
	Ed25519
)

func NewAlgorithm(name string) (error, Algorithm) {
	switch name {
	case "ed25519":
		return nil, Ed25519
	default:
		return ErrNotSupportAlgorithm, 0
	}
}

func (a Algorithm) getBackend() (error, AlgorithmI) {
	switch a {
	case Ed25519:
		return nil, &curver.Ed25519{}
	default:
		return ErrNotSupportAlgorithm, &curver.Ed25519{}
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

func (a Algorithm) GeneratePrivateKey() (error, []byte) {
	err, i := a.getBackend()
	if err != nil {
		return err, nil
	}
	return i.GeneratePrivateKey()
}

func (a Algorithm) CheckPrivateKey(privateKey []byte) error {
	err, i := a.getBackend()
	if err != nil {
		return err
	}
	return i.CheckPrivateKey(privateKey)
}

// GetPubkey will get the public key of the secret key
func (a Algorithm) GeneratePublicKey(privateKey []byte) (error, []byte) {
	err, i := a.getBackend()
	if err != nil {
		return err, nil
	}
	return i.GeneratePublicKey(privateKey)
}

func (a Algorithm) Sign(message []byte, privateKey []byte) (error, []byte) {
	err, i := a.getBackend()
	if err != nil {
		return err, nil
	}
	return i.Sign(message, privateKey)
}

func (a Algorithm) Verify(message []byte, publicKey []byte, sig []byte) (ret bool) {
	err, i := a.getBackend()
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
