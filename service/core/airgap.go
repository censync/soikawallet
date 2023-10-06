package core

import (
	"fmt"
	airgap "github.com/censync/go-airgap"
	"github.com/censync/soikawallet/api/dto"
)

func (s *Wallet) ProcessAirGapMessage(dto *dto.AirGapMessageDTO) (string, error) {
	airgapInstance := airgap.NewAirGap(airgap.VersionDefault, s.instanceId)

	message, err := airgapInstance.Unmarshal(dto.Data)

	if err != nil {
		return "", err
	}

	result := ""
	for _, operation := range message.Operations {
		result += fmt.Sprintf("op: %d  data_size(%d)", operation.OpCode, operation.Size)
	}
	return result, nil
}
