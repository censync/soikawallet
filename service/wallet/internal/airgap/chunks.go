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
	chunkHeaderOffset = 6                      // chunk_index(2) + chunks_count(2) + chunk_size(2)
	minMessageSize    = 35 + chunkHeaderOffset // protocol_version(1) + root_compressed_pub(33) + operation(1)
	baseChunkSize     = 192                    // best size for terminal
)

type Chunks struct {
	count  uint16
	size   uint16
	chunks [][]byte
}

func NewChunks(src []byte, chunkSize int) (*Chunks, error) {
	/*if len(src) < minMessageSize {
		return nil, errors.New("less than airgap message minimum size")
	}*/
	if chunkSize < minMessageSize {
		return nil, errors.New("min chunk size 32")
	}

	if chunkSize > 1<<16-chunkHeaderOffset {
		return nil, errors.New("max chunk size 65531")
	}

	chunkSize -= chunkHeaderOffset

	compressedData, err := compress(src)

	if err != nil {
		return nil, err
	}

	chunks := make([][]byte, 0)
	for iter := 0; iter < len(compressedData); iter += chunkSize {

		payloadSize := len(compressedData[iter:])

		chunk := make([]byte, 0)
		if payloadSize >= chunkSize {
			chunk = make([]byte, chunkSize)
			copy(chunk, compressedData[iter:iter+chunkSize])
		} else {
			chunk = make([]byte, payloadSize)
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
	size := len(f.chunks[index])
	chunk := make([]byte, f.size+chunkHeaderOffset)
	// chunk_index
	chunk[0] = byte(index)
	chunk[1] = byte(index >> 8)
	// chunk_count
	chunk[2] = byte(f.count)
	chunk[3] = byte(f.count >> 8)
	// chunk_size
	chunk[4] = byte(size)
	chunk[5] = byte(size >> 8)

	copy(chunk[chunkHeaderOffset:], f.chunks[index])

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

func (f *Chunks) ReadChunk(frame string) error {
	chunk, err := base64.StdEncoding.DecodeString(frame)

	if err != nil {
		return err
	}

	if f.count == 0 {
		f.count = uint16(chunk[2] | chunk[3]<<8)
		f.chunks = make([][]byte, f.count)
	}

	index := int16(chunk[0] | chunk[1]<<8)

	size := int16(chunk[4] | chunk[5]<<8)

	//fmt.Println(chunk[chunkHeaderOffset:])

	if f.chunks[index] == nil {
		f.chunks[index] = make([]byte, size)
		copy(f.chunks[index], chunk[chunkHeaderOffset:chunkHeaderOffset+size])
	}

	return nil
}

func (f *Chunks) Data() ([]byte, error) {
	data := make([]byte, 0)
	for i := 0; i < len(f.chunks); i++ {
		data = append(data, f.chunks[i]...)
	}
	result, err := uncompress(data)
	return result, err
}
