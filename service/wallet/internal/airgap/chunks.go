package airgap

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

const (
	chunkHeaderSize = 4   // chunk_index(2) + chunks_count(2)
	minMessageSize  = 18  // protocol_version(1) + uuid(16) + operation(1)
	baseChunkSize   = 192 // best size for terminal
)

type Chunks struct {
	count  uint16
	size   uint16
	chunks [][]byte
}

func NewChunks(src []byte, chunkSize int) (*Chunks, error) {
	if len(src) < minMessageSize {
		return nil, errors.New("less than airgap message minimum size")
	}
	if chunkSize < 32 {
		return nil, errors.New("min chunk size 32")
	}

	if chunkSize > 1<<16-4 {
		return nil, errors.New("max chunk size 65531")
	}

	chunkSize -= chunkHeaderSize

	compressedData, err := compress(src)

	if err != nil {
		return nil, err
	}

	chunks := make([][]byte, 0)
	for iter := 0; iter < len(compressedData); iter += chunkSize {
		chunk := make([]byte, chunkSize)
		if len(compressedData[iter:]) >= chunkSize {
			copy(chunk, compressedData[iter:iter+chunkSize])
		} else {
			copy(chunk, compressedData[iter:])
		}

		chunks = append(chunks, chunk)
	}

	return &Chunks{
		count:  uint16(len(chunks)),
		size:   uint16(chunkSize),
		chunks: chunks,
	}, nil
}

func compress(src []byte) ([]byte, error) {
	var buf bytes.Buffer
	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot compress data: %s", err.Error()))
	}

	_, err = zw.Write(src)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot write compressed data: %s", err.Error()))
	}

	if err = zw.Close(); err != nil {
		return nil, errors.New(fmt.Sprintf("cannot close writer: %s", err.Error()))
	}

	return buf.Bytes(), nil
}

func uncompress(src []byte) ([]byte, error) {
	reader := bytes.NewReader(src)

	zr, err := gzip.NewReader(reader)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot uncompress data: %s", err.Error()))
	}

	defer zr.Close()

	uncompressedBytes, err := io.ReadAll(zr)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot read uncompressed data: %s", err.Error()))
	}

	return uncompressedBytes, nil
}

func (f *Chunks) Chunk(index uint16) []byte {
	chunk := make([]byte, f.size)
	chunk[0] = byte(index)
	chunk[1] = byte(index >> 8)
	chunk[2] = byte(f.count)
	chunk[3] = byte(f.count >> 8)
	copy(chunk[chunkHeaderSize:], f.chunks[index])
	return chunk
}

func (f *Chunks) Bytes() [][]byte {
	var chunks [][]byte
	for index := uint16(0); index < f.count; index++ {
		chunks = append(chunks, f.Chunk(index))
	}
	return chunks
}

func (f *Chunks) ChunkBase64(index uint16) string {
	return base64.StdEncoding.EncodeToString(f.Chunk(index))
}

func (f *Chunks) ChunksBase64() []string {
	var chunksB64 []string
	for i := uint16(0); i < f.count; i++ {
		chunksB64 = append(chunksB64, base64.StdEncoding.EncodeToString(f.Chunk(i)))
	}
	return chunksB64
}

func (f *Chunks) Count() uint16 {
	return f.count
}
