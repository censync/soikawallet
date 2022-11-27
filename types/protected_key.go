package types

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"syscall"
)

type ProtectedKey struct {
	Key []byte
	Len int
}

func NewProtectedKey(key *ecdsa.PrivateKey) *ProtectedKey {
	pKey := &ProtectedKey{}
	pKey.Set(key)
	return pKey
}

func (m *ProtectedKey) lock() {
	if err := syscall.Mprotect(m.Key, syscall.PROT_NONE); err != nil {
		panic(err)
	}
}

func (m *ProtectedKey) unlock() {
	if err := syscall.Mprotect(m.Key, syscall.PROT_READ); err != nil {
		panic(err)
	}
}

func (m *ProtectedKey) Set(key *ecdsa.PrivateKey) {
	pageSize := syscall.Getpagesize()

	if len(crypto.FromECDSA(key)) > pageSize {
		panic("data larger than page size")
	}
	m.Key = make([]byte, pageSize)
	if err := syscall.Mprotect(m.Key, syscall.PROT_WRITE); err != nil {
		panic(err)
	}
	m.Len = copy(m.Key, crypto.FromECDSA(key))

	m.lock()

	key = nil
}

func (m *ProtectedKey) Get() *ecdsa.PrivateKey {

	m.unlock()
	defer m.lock()
	key, err := crypto.ToECDSA(m.Key[0:m.Len])

	if err != nil {
		panic(err)
	}

	return key
}

func (m *ProtectedKey) Free() {
	if m.Key != nil {
		if err := syscall.Mprotect(m.Key, syscall.PROT_WRITE|syscall.PROT_READ); err != nil {
			panic(err)
		}
	}
}
