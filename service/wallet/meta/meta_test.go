package meta

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestMeta_MarshalJSON(t *testing.T) {
	meta := InitMeta()

	meta.labelsAccount.Add("Account label 1")
	meta.labelsAccount.Add("Account label 2")
	meta.labelsAddress.Add("Address label 1")
	meta.labelsAddress.Add("Address label 2")

	str, err := json.Marshal(meta)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(str), err)

	meta2 := &Meta{}
	err = json.Unmarshal(str, &meta2)

	if err != nil {
		t.Fatal(err)
	}
}
