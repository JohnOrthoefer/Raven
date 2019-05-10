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
  "log"
  "path/filepath"
  "plugin"
	"../ravenTypes"
)

type CheckIMap map[string]ravenTypes.CheckInitType
type CheckFMap map[string]ravenTypes.CheckFuncType

var CheckInit CheckIMap
var CheckFunc CheckFMap

func init() {
  CheckFunc = make(CheckFMap)
	CheckInit = make(CheckIMap)

  plugins, err := filepath.Glob( "plugins/*.so")
  if err != nil {
    log.Fatal(err)
  }

  for i,v := range plugins {
    log.Printf("Plugin-%d: %s\n", i, v)

    plug,err := plugin.Open(v)
    if err != nil {
      log.Fatal( "Error openin %s", v)
    }

    n,err := plug.Lookup( "CheckName")
    if err != nil {
      log.Fatal( "No CheckName in %s", v)
    }
    checkName := n.(*string)

	  fInitCheck,err := plug.Lookup( "InitCheck")
    if err != nil {
      log.Fatal( "Error looking up InitCheck() in %s", v)
    }
	  fi,ok := fInitCheck.(func(ravenTypes.Kwargs) interface{})
    if !ok {
      log.Fatal( "No InitCheck() symbol in %s", v)
    }
    CheckInit[*checkName] = fi

	  fRunCheck,err := plug.Lookup( "RunCheck")
    if err != nil {
      log.Fatal( "Error looking up RunCheck() in %s", v)
    }
	  fc,ok := fRunCheck.(func(*ravenTypes.HostEntry, interface{}) *ravenTypes.ExitReturn)
    if !ok {
      log.Fatal( "No RunCheck() symbol in %s", v)
    }
    CheckFunc[*checkName] = fc
  }
}
