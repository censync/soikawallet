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

package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/censync/soikawallet/service/core/internal/config/version"
	"io"
	"net/http"
)

type ApiClient struct {
	Host string
	*http.Client
}

func (c *ApiClient) Do(method, path string, payload, dst interface{}) error {
	encodedPayload, err := json.Marshal(payload)

	if err != nil {
		return errors.New("cannot marshal request")
	}

	request, err := http.NewRequest(method, c.Host+path, bytes.NewBuffer(encodedPayload))

	if err != nil {
		return errors.New("cannot create request")
	}
	request.Header.Set("Content-Type", "application/json")

	request.Header.Set("User-Agent", fmt.Sprintf("Soika Wallet / %s", version.VERSION))
	request.Header.Set("TRON-PRO-API-KEY", "b57151c7-fa3e-4400-a382-2d53fc2ae5a0")
	resp, err := c.Client.Do(request)
	if err != nil {
		return errors.New("cannot process request: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("cannot execute request: code %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	return json.Unmarshal(data, &dst)
}
