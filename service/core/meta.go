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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

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
