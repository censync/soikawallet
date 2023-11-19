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

package subscriptions

import mhda "github.com/censync/go-mhda"

// EventsFilter map[chain_key]map[subscription_id][]Filter
type EventsFilter map[mhda.ChainKey]map[string][]*Filter

type FilterFunc func(arg ...interface{}) bool

type Filter struct {
	TaskId     string
	Conditions []FilterFunc
}
