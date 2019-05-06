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

package ravenChecks

import (
	"../ravenTypes"
)

type CheckIMap map[string]ravenTypes.CheckInitType
type CheckFMap map[string]ravenTypes.CheckFuncType

var CheckInit CheckIMap
var CheckFunc CheckFMap

func init() {
}

func registerHandler(n string,
	iptr ravenTypes.CheckInitType,
	fptr ravenTypes.CheckFuncType) {
	if CheckFunc == nil {
		CheckFunc = make(CheckFMap)
		CheckInit = make(CheckIMap)
	}
	CheckInit[n] = iptr
	CheckFunc[n] = fptr
}
