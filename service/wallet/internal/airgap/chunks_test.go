package airgap

import (
	"crypto/rand"
	"reflect"
	"testing"
)

func TestChunks_NewChunks(t *testing.T) {
	// without remainder
	/*chunksCount := uint16(3)

	payload := make([]byte, (defaultChunkSize-chunkHeaderOffset)*chunksCount)
	rand.Read(payload)

	chunksWithoutRemainder, err := NewChunks(payload, defaultChunkSize)

	if err != nil {
		t.Fatal(err)
	}

	if chunksWithoutRemainder.count != chunksCount {
		t.Fatal("data without remainder generation failed")
	}*/

	chunksCount := uint16(3)
	remainder := uint16(0)

	payload := make([]byte, defaultChunkSize*chunksCount-remainder)

	count, err := rand.Read(payload)

	if err != nil {
		t.Fatal("cannot read random")
	}

	t.Log("Readed random:", count)

	chunksWithRemainder, err := NewChunks(payload, defaultChunkSize)

	if err != nil {
		t.Fatal(err)
	}

	strChunks := chunksWithRemainder.SerializeB64()

	readedChunks := &chunks{}

	for i := 0; i < len(strChunks); i++ {
		err = readedChunks.ReadB64Chunk(strChunks[i])
		if err != nil {
			t.Fatal("cannot parse frame")
		}
	}

	result := make([]byte, 0)
	for i := 0; i < len(readedChunks.data); i++ {
		result = append(result, readedChunks.data[i]...)
	}
	uncompressedResult, err := uncompress(result)

	if err != nil {
		t.Fatal("cannot uncompress data", err)
	}

	if !reflect.DeepEqual(payload, uncompressedResult) {
		t.Fatal("mismatch marshalled data")
	}
}
