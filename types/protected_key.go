package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"reflect"
	"syscall"
)

type ProtectedKey struct {
	key []byte
	len int
}

func NewProtectedKey(key *ecdsa.PrivateKey) *ProtectedKey {
	pKey := &ProtectedKey{}
	pKey.set(key)
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

// set check for OSx https://github.com/apple-oss-distributions/Libsystem/tree/Libsystem-1336
func (pk *ProtectedKey) set(key *ecdsa.PrivateKey) {
	pageSize := syscall.Getpagesize()

	if len(crypto.FromECDSA(key)) > pageSize {
		panic("data larger than page size")
	}

	pk.key = make([]byte, pageSize)
	pk.len = copy(pk.key, crypto.FromECDSA(key))

	// pk.key size must allocate pageSize memory
	if err := syscall.Mprotect(pk.key, syscall.PROT_WRITE); err != nil {
		panic(err)
	}
	key = nil

	pk.lockMem()
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

func CheckMProtect() (ok bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
			err = errors.New("undefined error")
			return
		}
	}()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return false, errors.New("cannot generate key pair")
	}
	origKey := crypto.FromECDSA(key)

	pageSize := syscall.Getpagesize()

	bKey := make([]byte, pageSize)

	keyLen := copy(bKey, origKey)

	if keyLen > pageSize {
		return false, errors.New("data larger than page size")
	}

	if err = syscall.Mprotect(bKey, syscall.PROT_WRITE); err != nil {
		return false, errors.New("cannot write protected memory")
	}

	if err = syscall.Mprotect(bKey, syscall.PROT_NONE); err != nil {
		return false, errors.New("cannot lock memory")
	}

	if err = syscall.Mprotect(bKey, syscall.PROT_READ); err != nil {
		return false, errors.New("cannot unlock readable memory")
	}

	ok = reflect.DeepEqual(bKey[0:keyLen], origKey)

	if !ok {
		return false, errors.New("invalid memory data")
	}

	if err = syscall.Mprotect(bKey, syscall.PROT_WRITE|syscall.PROT_READ); err != nil {
		return false, errors.New("cannot flush protected memory")
	}

	return true, nil
}
