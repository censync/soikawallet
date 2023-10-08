// Copyright 2023 The soikawallet Authors
// This file is part of the soikawallet library.
//
// The soikawallet library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The soikawallet library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the  soikawallet library. If not, see <http://www.gnu.org/licenses/>.

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
