package airgap

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	ecies "github.com/ecies/go/v2"
	"testing"
)

const (
	opCodeTest1 = 1
	opCodeTest2 = 1000
	opCodeTest3 = 10000
)

type DummyEncryptorDecryptor struct {
	privKey *ecies.PrivateKey
	pubKey  *ecies.PublicKey
}

func NewDummyEncryptorDecryptor() *DummyEncryptorDecryptor {
	privKey, _ := ecies.GenerateKey()
	return &DummyEncryptorDecryptor{
		privKey: privKey,
		pubKey:  privKey.PublicKey,
	}
}

func (ed *DummyEncryptorDecryptor) Encrypt(data []byte) ([]byte, error) {
	return ecies.Encrypt(ed.pubKey, data)
}

func (ed *DummyEncryptorDecryptor) Decrypt(data []byte) ([]byte, error) {
	return ecies.Decrypt(ed.privKey, data)
}

func TestAirGap_CreateMessage(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal("cannot generate private key")
	}

	pubKeySerialized := elliptic.MarshalCompressed(elliptic.P256(), privKey.X, privKey.Y)

	airGap := NewAirGap(VersionDefault, pubKeySerialized)

	// ed := NewDummyEncryptorDecryptor()

	// airGap.SetEncryptorDecryptor(ed)

	opMessage := airGap.CreateMessage().
		AddOperation(opCodeTest1, []byte("secret message 1")).
		AddOperation(opCodeTest2, []byte("secret message 2")).
		AddOperation(opCodeTest3, []byte("secret message 3"))

	serializedChunks, err := opMessage.MarshalB64Chunks()
	serializedChunks2, err := opMessage.MarshalB64Chunks()
	serializedChunks3, err := opMessage.MarshalB64Chunks()
	t.Log(serializedChunks)
	t.Log(serializedChunks2)
	t.Log(serializedChunks3)
	if err != nil {
		t.Fatal("cannot serialize chunks")
	}

	// t.Log(serializedChunks)

	unserializedChunks := &chunks{}

	for i := range serializedChunks {
		err = unserializedChunks.ReadB64Chunk(serializedChunks[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	chunksData := unserializedChunks.Data()

	t.Log(chunksData)
	unserializedMessage, err := airGap.Unmarshal(chunksData)

	if err != nil {
		t.Log(err)
		t.Fatal("cannot unserialize chunks")
	}

	t.Log(unserializedMessage)
}
