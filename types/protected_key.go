package types

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"syscall"
)

type ProtectedKey struct {
	key []byte
	len int
}

func NewProtectedKey(key *ecdsa.PrivateKey) *ProtectedKey {
	pKey := &ProtectedKey{}
	pKey.Set(key)
	return pKey
}

func (pk *ProtectedKey) lockMem() {
	if err := syscall.Mprotect(pk.key, syscall.PROT_NONE); err != nil {
		panic(err)
	}
}

func (pk *ProtectedKey) unlockMem() {
	if err := syscall.Mprotect(pk.key, syscall.PROT_READ); err != nil {
		panic(err)
	}
}

func (pk *ProtectedKey) Set(key *ecdsa.PrivateKey) {
	pageSize := syscall.Getpagesize()

	if len(crypto.FromECDSA(key)) > pageSize {
		panic("data larger than page size")
	}
	pk.key = make([]byte, pageSize)
	if err := syscall.Mprotect(pk.key, syscall.PROT_WRITE); err != nil {
		panic(err)
	}
	pk.len = copy(pk.key, crypto.FromECDSA(key))

	pk.lockMem()

	key = nil
}

func (pk *ProtectedKey) Get() *ecdsa.PrivateKey {

	pk.unlockMem()
	defer pk.lockMem()
	key, err := crypto.ToECDSA(pk.key[0:pk.len])

	if err != nil {
		panic(err)
	}

	return key
}

func (pk *ProtectedKey) Free() {
	if pk.key != nil {
		if err := syscall.Mprotect(pk.key, syscall.PROT_WRITE|syscall.PROT_READ); err != nil {
			panic(err)
		}
	}
}
