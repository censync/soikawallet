package wallet

import (
	"errors"
	"github.com/censync/soikawallet/api/dto"
	"github.com/censync/soikawallet/types"
	"strings"
)

var (
	errMetaLabelUnknownType = errors.New("unknown label type")
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

	if dto.LabelType != types.AccountLabel && dto.LabelType != types.AddressLabel {
		return 0, errors.New("unknown label type")
	}

	switch dto.LabelType {
	case types.AccountLabel:
		return s.meta.AddAccountLabel(dto.Title)
	case types.AddressLabel:
		return s.meta.AddAddressLabel(dto.Title)
	default:
		return 0, errMetaLabelUnknownType
	}
}

func (s *Wallet) RemoveLabel(dto *dto.RemoveLabelDTO) error {
	switch dto.LabelType {
	case types.AccountLabel:
		return s.meta.RemoveAccountLabel(dto.Index)
	case types.AddressLabel:
		return s.meta.RemoveAddressLabel(dto.Index)
	default:
		return errMetaLabelUnknownType
	}
}

func (s *Wallet) SetLabelLink(dto *dto.SetLabelLinkDTO) error {
	switch dto.LabelType {
	case types.AccountLabel:
		return s.meta.SetAccountLabelLink(dto.Path, dto.Index)
	case types.AddressLabel:
		return s.meta.SetAddressLabelLink(dto.Path, dto.Index)
	default:
		return errMetaLabelUnknownType
	}
}

func (s *Wallet) RemoveLabelLink(dto *dto.RemoveLabelLinkDTO) error {
	switch dto.LabelType {
	case types.AccountLabel:
		return s.meta.RemoveAccountLabelLink(dto.Path)
	case types.AddressLabel:
		return s.meta.RemoveAddressLabelLink(dto.Path)
	default:
		return errMetaLabelUnknownType
	}
}
