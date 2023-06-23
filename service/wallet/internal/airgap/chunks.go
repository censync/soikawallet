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
	defaultChunkSize  = 192                    // best size for terminal
	minChunkSize      = 34 + chunkHeaderOffset // protocol_version(1) + root_compressed_pub(33)
	// maxChunksCount = maxChunksCount
	// maxPayloadSize    = defaultChunkSize * (2<<15 - 1 - chunkHeaderOffset) // ~ 12.58Mb
)

type chunks struct {
	count uint16
	size  uint16
	data  [][]byte
}

func NewChunks(src []byte, chunkSize int) (*chunks, error) {
	/*if len(src) < minChunkSize {
		return nil, errors.New("less than airgap message minimum size")
	}*/
	if chunkSize < minChunkSize {
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

	data := make([][]byte, 0)
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

		data = append(data, chunk)
	}

	return &chunks{
		count: uint16(len(data)),
		size:  uint16(chunkSize),
		data:  data,
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

func (f *chunks) getChunkWithHeader(index uint16) []byte {
	size := len(f.data[index])
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

	copy(chunk[chunkHeaderOffset:], f.data[index])

	return chunk
}

func (f *chunks) Data() []byte {
	var result []byte
	for index := uint16(0); index < f.count; index++ {
		result = append(result, f.data[index]...)
	}
	result, _ = uncompress(result)
	return result
}

// SerializeB64 represents data frames to strings array, ready for generate QR code animation frames
func (f *chunks) SerializeB64() []string {
	var chunksB64 []string
	for i := uint16(0); i < f.count; i++ {
		chunksB64 = append(chunksB64, base64.StdEncoding.EncodeToString(f.getChunkWithHeader(i)))
	}
	return chunksB64
}

func (f *chunks) Count() uint16 {
	return f.count
}

func (f *chunks) ReadB64Chunk(frame string) error {
	chunk, err := base64.StdEncoding.DecodeString(frame)

	if err != nil {
		return err
	}

	if f.count == 0 {
		f.count = uint16(chunk[2]) | uint16(chunk[3])<<8
		f.data = make([][]byte, f.count)
	}

	index := uint16(chunk[0]) | uint16(chunk[1])<<8

	size := uint16(chunk[4]) | uint16(chunk[5])<<8

	if f.data[index] == nil {
		f.data[index] = make([]byte, size)
		copy(f.data[index], chunk[chunkHeaderOffset:chunkHeaderOffset+size])
	}

	return nil
}

func (f *chunks) IsFilled() bool {
	return len(f.data) == int(f.count)
}
