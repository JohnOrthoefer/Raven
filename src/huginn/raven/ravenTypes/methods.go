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
	"regexp"
	"strconv"
	"strings"
)

var reSpaces *regexp.Regexp

func init() {
	reSpaces = regexp.MustCompile(`\s+`)
}

func (kw Kwargs) GetKwargStr(s string, d string) string {
	s = strings.ToLower(s)
	if rtn, ok := kw[s]; ok {
		return rtn
	}
	return d
}

func (kw Kwargs) GetKwargStrTrim(s string, d string) string {
	return strings.TrimSpace(kw.GetKwargStr(s, d))
}

func (kw Kwargs) GetKwargStrA(s string, d []string) []string {
	s = strings.ToLower(s)
	if t, ok := kw[s]; ok {
		return reSpaces.Split(strings.TrimSpace(t), -1)
	}
	return d
}

func (kw Kwargs) GetKwargBool(s string, d bool) bool {
	s = strings.ToLower(s)
	if t, ok := kw[s]; ok {
		t = strings.ToLower(t)
		for _, i := range []string{"true", "t", "yes", "y"} {
			if i == t {
				return true
			}
		}
		return false
	}
	return d
}

func (kw Kwargs) GetKwargFloat(s string, d float64) float64 {
	s = strings.ToLower(s)
	if rtn, ok := kw[s]; ok {
		if fl, err := strconv.ParseFloat(rtn, 64); err == nil {
			return fl
		}
	}
	return d
}

func (kw Kwargs) GetKwargInt(s string, d int64) int64 {
	s = strings.ToLower(s)
	if rtn, ok := kw[s]; ok {
		if num, err := strconv.ParseInt(rtn, 10, 0); err == nil {
			return num
		}
	}
	return d
}
