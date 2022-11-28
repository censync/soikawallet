package airgap

import (
	"crypto/rand"
	"testing"
)

func TestAirGap_CreateMessage(t *testing.T) {
	var metaStr = `{"meta":{"v":1,"labels":{"1":{"1":"Account label 1","2":"Account label 2"},"2":{"1":"addr label 1","2":"addr label 2"}}},"addresses":["m/44'/60'/0'/0/9'","m/44'/60'/0'/0/0'","m/44'/60'/0'/0/1'","m/44'/60'/0'/0/4'","m/44'/60'/0'/0/6'","m/44'/60'/0'/0/7'","m/44'/60'/0'/0/8'","m/44'/60'/0'/0/2'","m/44'/60'/0'/0/3'","m/44'/60'/0'/0/5'","m/44'/195'/0'/0/2'","m/44'/195'/0'/0/5'","m/44'/195'/0'/0/6'","m/44'/195'/0'/0/7'","m/44'/195'/0'/0/0'","m/44'/195'/0'/0/1'","m/44'/195'/0'/0/3'","m/44'/195'/0'/0/4'","m/44'/195'/0'/0/8'","m/44'/195'/0'/0/9'"]}`

	airgapInstance := Create()

	message := airgapInstance.CreateMessage(OperationExportMeta, []byte(metaStr))

	t.Log(len(message.Bytes()))

	chunks, err := NewChunks(message.Bytes(), baseChunkSize)

	if err != nil {
		t.Fatal(err)
	}

	for i := uint16(0); i < chunks.Count(); i++ {
		t.Log(chunks.ChunkBase64(i))
	}
}

func TestChunks_NewChunks(t *testing.T) {
	// without remainder
	chunksCount := uint16(5)

	payload := make([]byte, baseChunkSize*chunksCount)
	rand.Read(payload)

	chunksWithoutRemainder, err := NewChunks(payload, baseChunkSize)

	if err != nil {
		t.Fatal(err)
	}

	if chunksWithoutRemainder.count != chunksCount {
		t.Fatal("chunks without remainder generation failed")
	}

	remainder := uint16(7)

	payload = make([]byte, baseChunkSize*chunksCount-remainder)
	rand.Read(payload)

	chunksWithRemainder, err := NewChunks(payload, baseChunkSize)

	if err != nil {
		t.Fatal(err)
	}

	if chunksWithRemainder.count != chunksCount {
		t.Fatal("chunks with remainder generation failed")
	}

}
