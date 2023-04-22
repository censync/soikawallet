package wallet

import (
	"encoding/base64"
	"github.com/censync/soikawallet/service/wallet/internal/airgap"
	"github.com/censync/soikawallet/types"
	ecies2 "github.com/ecies/go/v2"
)

func AirGapInitMessage() string {
	a := airgap.Create()
	airGapMessage := a.CreateMessage().AddOperation(types.OpInitBootstrap, nil)
	b64 := base64.StdEncoding.EncodeToString(airGapMessage.Bytes())
	return b64
}

func GetECIESPub() []byte {
	k, err := ecies2.GenerateKey()
	if err != nil {
		panic(err)
	}
	return k.PublicKey.Bytes(true)
}

/*
func EncryptOAEP(hash hash.Hash, random io.Reader, public *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := public.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, public, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

func DecryptOAEP(hash hash.Hash, random io.Reader, private *rsa.PrivateKey, msg []byte, label []byte) ([]byte, error) {
	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, private, msg[start:finish], label)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
*/
