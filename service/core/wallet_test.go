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

package core

import (
	"github.com/censync/soikawallet/api/dto"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testMnemonic = `ceiling code ginger fabric transfer gallery sort deputy engine blur believe sunny divert nephew brain tired result husband upper clock auction ritual correct inhale`

func TestWalletService_Init(t *testing.T) {
	service := &Wallet{}
	_, err := service.Init(&dto.InitWalletDTO{
		Mnemonic:   testMnemonic,
		Passphrase: ``,
	})
	assert.Nil(t, err)
	//priv, _ := service.rootKey.ECPrivKey()
	//pub, _ := service.rootKey.ECPubKey()
	//t.Log(priv.key.String())
	//t.Log(pub.Y(), pub.Y(), err)
}
