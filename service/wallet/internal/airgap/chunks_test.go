package airgap

import (
	"crypto/rand"
	"github.com/censync/soikawallet/types"
	"reflect"
	"testing"
)

func TestAirGap_CreateMessage(t *testing.T) {
	var metaStr = `{"meta":{"v":1,"labels":{"1":{"1":"Account label 1","2":"Account label 2"},"2":{"1":"addr label 1","2":"addr label 2"}}},"addresses":["m/44'/60'/0'/0/9'","m/44'/60'/0'/0/0'","m/44'/60'/0'/0/1'","m/44'/60'/0'/0/4'","m/44'/60'/0'/0/6'","m/44'/60'/0'/0/7'","m/44'/60'/0'/0/8'","m/44'/60'/0'/0/2'","m/44'/60'/0'/0/3'","m/44'/60'/0'/0/5'","m/44'/195'/0'/0/2'","m/44'/195'/0'/0/5'","m/44'/195'/0'/0/6'","m/44'/195'/0'/0/7'","m/44'/195'/0'/0/0'","m/44'/195'/0'/0/1'","m/44'/195'/0'/0/3'","m/44'/195'/0'/0/4'","m/44'/195'/0'/0/8'","m/44'/195'/0'/0/9'"]}`

	airgapInstance := Create()

	message := airgapInstance.CreateMessage().AddOperation(types.OpMetaAirGap, []byte(metaStr))

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
	/*chunksCount := uint16(3)

	payload := make([]byte, (baseChunkSize-chunkHeaderOffset)*chunksCount)
	rand.Read(payload)

	chunksWithoutRemainder, err := NewChunks(payload, baseChunkSize)

	if err != nil {
		t.Fatal(err)
	}

	if chunksWithoutRemainder.count != chunksCount {
		t.Fatal("chunks without remainder generation failed")
	}*/

	chunksCount := uint16(3)
	remainder := uint16(0)

	payload := make([]byte, baseChunkSize*chunksCount-remainder)

	count, err := rand.Read(payload)

	if err != nil {
		t.Fatal("cannot read random")
	}

	t.Log("Readed random:", count)

	chunksWithRemainder, err := NewChunks(payload, baseChunkSize)

	if err != nil {
		t.Fatal(err)
	}

	strChunks := chunksWithRemainder.ChunksBase64()

	readedChunks := &Chunks{}

	for i := 0; i < len(strChunks); i++ {
		err = readedChunks.ReadChunk(strChunks[i])
		if err != nil {
			t.Fatal("cannot parse frame")
		}
	}

	result, err := readedChunks.Data()

	if err != nil {
		t.Fatal("cannot uncompress chunks", err)
	}

	if !reflect.DeepEqual(payload, result) {
		t.Fatal("mismatch marshalled data")
	}
}
