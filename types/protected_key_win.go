//go:build windows

// Syscall is not implemented on windows
package types

import (
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	memoryProtectionAvailable = false
	memoryProtectionError     error
)

type ProtectedKey struct {
	key []byte
}

func init() {
	memoryProtectionAvailable = false
	memoryProtectionError = errors.New("not implemented for windows")
}

func NewProtectedKey(key *ecdsa.PrivateKey) *ProtectedKey {
	return &ProtectedKey{
		key: crypto.FromECDSA(key),
	}
}

func (pk *ProtectedKey) Get() *ecdsa.PrivateKey {
	key, err := crypto.ToECDSA(pk.key)
	if err != nil {
		panic(err)
	}
	return key
}

func (pk *ProtectedKey) Free() {
	if pk.key != nil {
		pk.key = nil
	}
}
