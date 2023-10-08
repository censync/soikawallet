// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

//go:build !windows

package protected_key

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"reflect"
	"syscall"
)

var (
	memoryProtectionAvailable = false
	memoryProtectionError     error
)

type ProtectedKey struct {
	key []byte
	len int
}

func init() {
	defer func() {
		if r := recover(); r != nil {
			memoryProtectionError = errors.New("undefined error")
			return
		}
	}()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		memoryProtectionError = errors.New("cannot generate key pair")
		return
	}
	origKey := crypto.FromECDSA(key)

	pageSize := syscall.Getpagesize()

	bKey := make([]byte, pageSize)

	keyLen := copy(bKey, origKey)

	if keyLen > pageSize {
		memoryProtectionError = errors.New("data larger than pages size")
		return
	}

	if err = syscall.Mprotect(bKey, syscall.PROT_WRITE); err != nil {
		memoryProtectionError = errors.New("cannot write protected memory")
		return
	}

	if err = syscall.Mprotect(bKey, syscall.PROT_NONE); err != nil {
		memoryProtectionError = errors.New("cannot lock memory")
		return
	}

	if err = syscall.Mprotect(bKey, syscall.PROT_READ); err != nil {
		memoryProtectionError = errors.New("cannot unlock readable memory")
		return
	}

	ok := reflect.DeepEqual(bKey[0:keyLen], origKey)

	if !ok {
		memoryProtectionError = errors.New("invalid memory data")
		return
	}

	if err = syscall.Mprotect(bKey, syscall.PROT_WRITE|syscall.PROT_READ); err != nil {
		memoryProtectionError = errors.New("cannot flush protected memory")
		return
	}

	origKey = nil
	bKey = nil
	memoryProtectionAvailable = true
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
// TODO: Add mapping per pages size
func (pk *ProtectedKey) set(key *ecdsa.PrivateKey) {
	if memoryProtectionAvailable {
		pageSize := syscall.Getpagesize()

		if len(crypto.FromECDSA(key)) > pageSize {
			panic("data larger than pages size")
		}

		pk.key = make([]byte, pageSize)
		pk.len = copy(pk.key, crypto.FromECDSA(key))

		// pk.key size must multiple to pageSize memory
		// Known issue on macOS >= Catalina
		// https://stackoverflow.com/questions/60654834/using-mprotect-to-make-text-segment-writable-on-macos
		// l64 implementation: https://github.com/apple-opensource/ld64/blob/fd3feabb0a1eb18ab5d7910f3c3a5eed99cef6ab/src/ld/Options.cpp#L374-L379
		if err := syscall.Mprotect(pk.key, syscall.PROT_WRITE); err != nil {
			panic(err)
		}
		pk.lockMem()
	} else {
		pk.key = make([]byte, len(crypto.FromECDSA(key)))
		pk.len = copy(pk.key, crypto.FromECDSA(key))
	}
	key = nil
}

func (pk *ProtectedKey) Get() *ecdsa.PrivateKey {
	if memoryProtectionAvailable {
		pk.unlockMem()
		defer pk.lockMem()

		key, err := crypto.ToECDSA(pk.key[0:pk.len])

		if err != nil {
			panic(err)
		}

		return key
	} else {
		key, err := crypto.ToECDSA(pk.key)
		if err != nil {
			panic(err)
		}
		return key
	}
}

func (pk *ProtectedKey) Free() {
	if pk.key != nil {
		if memoryProtectionAvailable {
			if err := syscall.Mprotect(pk.key, syscall.PROT_WRITE|syscall.PROT_READ); err != nil {
				panic(err)
			}
		}
		pk.key = nil
	}
}

func IsMemoryProtected() (bool, error) {
	return memoryProtectionAvailable, memoryProtectionError
}
