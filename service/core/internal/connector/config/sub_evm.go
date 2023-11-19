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

package config

import (
	mhda "github.com/censync/go-mhda"
	sub "github.com/censync/soikawallet/service/core/internal/connector/subscriptions"
	log "github.com/sirupsen/logrus"
)

var Subs = map[mhda.ChainKey][]*sub.Subscription{
	"ethereum": SubsEVMs,
	"polygon":  SubsEVMs,
}

var headsFunc = func(event interface{}) {
	_, ok := event.(map[string]interface{})
	log.Infof("Got subscription response: %v", event)
	if !ok {
		log.Warnf("Cannot unmarshal sub response")

	}
	/*pHash := subRespJ["parentHash"].(string)
		currentHash := subRespJ["hash"].(string)

		blockNumStr := subRespJ["number"].(string)

		blockNumStr = strings.TrimLeft(blockNumStr, "0x")
		currBlockNum, err := strconv.ParseUint(blockNumStr, 16, 64)

		if err != nil {
			log.Warn("Cannot convert block num", blockNumStr, err)
		}

		if parentHash != "" {
			log.Infof("Got subscription response %s %d block %v: %v", chainKey, index, currBlockNum, event)
			if parentHash != pHash {
				log.Error("PARENT HASH MISTMATCH")
				log.Errorf("Parent: %d, current: %d", parentBlockNum, currBlockNum)
				log.Errorf("%s vs %s", parentHash, pHash)
			} else {
				parentBlockNum = currBlockNum
				parentHash = currentHash
			}
		} else {
			parentBlockNum = currBlockNum
			parentHash = currentHash
		}
	}*/
}

var SubsEVMs = []*sub.Subscription{
	{
		Name:    "base_height",
		Impl:    sub.Socket,
		SubType: sub.Latest,
		Status:  sub.Idle,
		Method:  "eth_subscribe",
		Params:  []interface{}{"newHeads"},
		Func:    headsFunc,
	},
	/*{
		Name:    "base_syncing",
		Impl:    sub.Socket,
		SubType: sub.Syncing,
		Status:  sub.Idle,
		Method:  "eth_subscribe",
		Params:  []interface{}{"syncing"},
	},
	{
		Impl:    sub.Dummy,
		SubType: sub.Latest,
		Status:  sub.Idle,
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{"finalized", false},
	},*/
}
