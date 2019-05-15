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

package raven

import (
	"./ravenLog"
	"./ravenTypes"
	"fmt"
  "log"
  "github.com/go-ini/ini"
)

func makeMap(s *ini.Section) ravenTypes.Kwargs {
	r := make(ravenTypes.Kwargs)
	for _, key := range s.Keys() {
		n := key.Name()
		v := key.Value()
		r[n] = v
	}
	return r
}

func ReadConfig(iniFile string) {
	cfg, err := ini.Load(iniFile)
  if err != nil {
    ravenLog.SendError( 10, "Configuration", fmt.Sprintf("Can not read file %s", iniFile))
    log.Fatal(err)
  }
	for _, section := range cfg.Sections() {
		secName := section.Name()
		ravenLog.SendError(10, "Configuration", fmt.Sprintf("Parsing Section %s", secName))
		keyVal := makeMap(section)
		AddEntry(secName, keyVal)
	}
}
