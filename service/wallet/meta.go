package wallet

import (
	"errors"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/service/wallet/meta"
	"strings"
)

func (s *Wallet) GetAccountLabels() map[uint32]string {
	return s.meta.AccountLabels()
}

func (s *Wallet) GetAddressLabels() map[uint32]string {
	return s.meta.AddressLabels()
}

func (s *Wallet) AddLabel(dto *dto.AddLabelDTO) (uint32, error) {
	dto.Title = strings.TrimSpace(dto.Title)

	if len(dto.Title) < 2 {
		return 0, errors.New("minimum 2 chars")
	}

	if len(dto.Title) > 50 {
		return 0, errors.New("maximum 20 chars")
	}

	if meta.LabelType(dto.LabelType) != meta.AccountLabel && meta.LabelType(dto.LabelType) != meta.AddressLabel {
		return 0, errors.New("unknown label type")
	}

	switch meta.LabelType(dto.LabelType) {
	case meta.AccountLabel:
		return s.meta.AddAccountLabel(dto.Title)
	case meta.AddressLabel:
		return s.meta.AddAddressLabel(dto.Title)
	default:
		return 0, errors.New("unknown label type")
	}
}

func (s *Wallet) RemoveLabel(dto *dto.RemoveLabelDTO) error {
	switch meta.LabelType(dto.LabelType) {
	case meta.AccountLabel:
		return s.meta.RemoveAccountLabel(dto.Index)
	case meta.AddressLabel:
		return s.meta.RemoveAddressLabel(dto.Index)
	default:
		return errors.New("unknown label type")
	}
}
