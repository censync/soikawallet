package types

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const (
	testRPCTitle     = "My test RPC"
	testRPCEndpoint  = "https://rpc.example.com/testnet"
	testRPCIsDefault = true
)

func TestRPC_MarshalJSON(t *testing.T) {
	rpc := NewRPC(
		testRPCTitle,
		testRPCEndpoint,
		testRPCIsDefault,
	)

	strJSON, err := json.Marshal(rpc)
	assert.Nil(t, err)

	assert.Equal(t, []byte(fmt.Sprintf(
		`["%s","%s","%s"]`,
		testRPCTitle,
		testRPCEndpoint,
		strconv.FormatBool(testRPCIsDefault)),
	),
		strJSON,
	)
}
