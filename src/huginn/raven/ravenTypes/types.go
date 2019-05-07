/*
   Raven Network Discovery and Monitoring
   Copyright (C) 2019 John{at}Orthoefer{dot}org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.

*/

package ravenTypes

import (
	"time"
)

type Kwargs map[string]string
type CheckInitType func(Kwargs) interface{}
type CheckFuncType func(*HostEntry, interface{}) *ExitReturn

type HostEntry struct {
	IPv4        string
	Hostname    string
	DisplayName string
	Group       string
}

type ExitReturn struct {
	Exit int
	Text string
	Perf string
	Long string
}

type CheckEntry struct {
	DisplayName string
	CheckF      CheckFuncType
	Options     interface{}
	CheckN      string
	Interval    [4]time.Duration
	Hosts       []string
	Threshold   int
}

// the basic entry for scheduling a check against a host
type StatusEntry struct {
	Check   *CheckEntry
	Host    *HostEntry
	Queued  bool
	CurExit int
	Count   int
	Next    time.Time
	Last    time.Time
	Change  time.Time
	Return  *ExitReturn
}
